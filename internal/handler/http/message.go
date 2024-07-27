package http

import (
	"chat-backend/internal/apperror"
	"chat-backend/internal/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	ID     uint64
	Socket *websocket.Conn
	Send   chan []byte
}

var clients = make(map[uint64]*Client)
var broadcast = make(chan model.ViewMessage)

type sendMessageRequest struct {
	UserID  uint64 `json:"user_id"`
	Content string `json:"content"`
}

func (h *Handler) sendMessage(c echo.Context) error {
	var req sendMessageRequest

	sender, err := getUserFromContext(c)
	if err != nil {
		return h.handleAppError(c, err)
	}

	err = c.Bind(&req)
	if err != nil {
		return h.handleAppError(c, apperror.ErrInvalidRequestBody)
	}

	createdMessage, err := h.services.Message.Create(c.Request().Context(), model.CreateMessageInput{
		SenderID:    sender.ID,
		RecipientID: req.UserID,
		Content:     req.Content,
	})
	if err != nil {
		return h.handleAppError(c, err)
	}
	broadcast <- createdMessage

	return ProcessDataHTTPResponse(c, http.StatusCreated, nil)
}

type getFriendMessagesResponse struct {
	Messages   []model.ViewMessage  `json:"messages"`
	Pagination model.FullPagination `json:"pagination"`
}

func (h *Handler) getFriendMessages(c echo.Context) error {
	user, err := getUserFromContext(c)
	if err != nil {
		return h.handleAppError(c, err)
	}

	pagination, err := h.getPaginationFromContext(c)
	if err != nil {
		return h.handleAppError(c, err)
	}

	id := c.Param("id")
	friendID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return h.handleAppError(c, apperror.GetErrInvalidParam("id"))
	}

	messages, fullPagination, err := h.services.Message.GetAllWithFriend(c.Request().Context(), pagination, user.ID, friendID)
	if err != nil {
		return h.handleAppError(c, err)
	}

	return ProcessDataHTTPResponse(c, http.StatusOK, getFriendMessagesResponse{
		Messages:   messages,
		Pagination: fullPagination,
	})

}

func (h *Handler) handleMessages(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-broadcast:
			messageJSON, _ := json.Marshal(msg)
			if _, ok := clients[msg.RecipientID]; ok {
				clients[msg.RecipientID].Send <- messageJSON
			}

			if _, ok := clients[msg.SenderID]; ok {
				clients[msg.SenderID].Send <- messageJSON
			}
		}
	}
}

func (c *Client) readPump(ctx context.Context) {
	defer func() {
		c.Socket.Close()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			fmt.Printf("read message\n")
			_, _, err := c.Socket.ReadMessage()
			if err != nil {
				fmt.Printf("read message error: %s\n", err)
				c.Socket.Close()
				return
			}
		}
	}
}

func (c *Client) writePump(ctx context.Context) {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.Socket.Close()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case message, ok := <-c.Send:
			if !ok {
				c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := c.Socket.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				return
			}
		case <-ticker.C:
			err := c.Socket.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				return
			}
		}
	}
}

func (h *Handler) listenMessageUpdates(c echo.Context) error {
	user, err := getUserFromContext(c)
	if err != nil {
		return h.handleAppError(c, err)
	}

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer conn.Close()
	h.logger.Debugf("new user connected: %s", user.Username)

	client := &Client{ID: user.ID, Socket: conn, Send: make(chan []byte)}

	clients[user.ID] = client

	go client.writePump(c.Request().Context())
	go client.readPump(c.Request().Context())
	h.logger.Debugf("new user pumps started: %s", user.Username)
	<-c.Request().Context().Done()
	delete(clients, user.ID)

	return nil
}
