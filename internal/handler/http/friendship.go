package http

import (
	"chat-backend/internal/model"
	"github.com/labstack/echo/v4"
	"net/http"
)

type friendshipsResponse struct {
	Friends    []model.ViewUser     `json:"friends"`
	Pagination model.FullPagination `json:"pagination"`
}

func (h *Handler) getMyFriendships(c echo.Context) error {
	user, err := getUserFromContext(c)
	if err != nil {
		return h.handleAppError(c, err)
	}

	pagination, err := h.getPaginationFromContext(c)
	if err != nil {
		return h.handleAppError(c, err)
	}

	friendships, fullPagination, err := h.services.Friendship.GetByUserID(c.Request().Context(), pagination, user.ID)
	if err != nil {
		return h.handleAppError(c, err)
	}
	return ProcessDataHTTPResponse(c, http.StatusOK, friendshipsResponse{
		Friends:    friendships,
		Pagination: fullPagination,
	})
}

type addFriendRequest struct {
	UserID uint64 `json:"user_id"`
}

func (h *Handler) AddFriend(c echo.Context) error {
	var req addFriendRequest

	user, err := getUserFromContext(c)
	if err != nil {
		return h.handleAppError(c, err)
	}

	err = c.Bind(&req)
	if err != nil {
		return h.handleAppError(c, err)
	}

	err = h.services.Friendship.Create(c.Request().Context(), model.CreateFriendshipDTO{
		UserID:   user.ID,
		FriendID: req.UserID,
	})
	if err != nil {
		return h.handleAppError(c, err)
	}

	return ProcessDataHTTPResponse(c, http.StatusCreated, nil)
}
