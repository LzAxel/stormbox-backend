package http

import (
	"chat-backend/internal/apperror"
	"chat-backend/internal/model"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

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

	_, err = h.services.Message.Create(c.Request().Context(), model.CreateMessageInput{
		SenderID:    sender.ID,
		RecipientID: req.UserID,
		Content:     req.Content,
	})
	if err != nil {
		return h.handleAppError(c, err)
	}

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
