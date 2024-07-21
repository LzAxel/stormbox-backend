package http

import (
	"chat-backend/internal/apperror"
	"chat-backend/internal/model"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type getAllUsersResponse struct {
	Users      []model.ViewUser     `json:"users"`
	Pagination model.FullPagination `json:"pagination"`
}

func (h *Handler) getAllUsers(ctx echo.Context) error {
	reqPagination, err := h.getPaginationFromContext(ctx)
	if err != nil {
		return h.handleAppError(ctx, err)
	}

	users, pagination, err := h.services.User.GetAll(ctx.Request().Context(), reqPagination)
	if err != nil {
		return h.handleAppError(ctx, err)
	}

	return ProcessDataHTTPResponse(ctx, http.StatusOK, getAllUsersResponse{Users: users, Pagination: pagination})
}

type getUserResponse struct {
	User model.ViewUser `json:"user"`
}

func (h *Handler) getUserByID(ctx echo.Context) error {
	id := ctx.Param("id")
	userID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return apperror.New(apperror.ErrorTypeBadRequest, "invalid user id", err)
	}

	user, err := h.services.User.GetByID(ctx.Request().Context(), userID)
	if err != nil {
		return h.handleAppError(ctx, err)
	}

	return ProcessDataHTTPResponse(ctx, http.StatusOK, getUserResponse{User: user})
}

func (h *Handler) getSelfUser(ctx echo.Context) error {
	ctxuser, err := getUserFromContext(ctx)
	if err != nil {
		return h.handleAppError(ctx, err)
	}

	user, err := h.services.User.GetByID(ctx.Request().Context(), ctxuser.ID)
	if err != nil {
		return h.handleAppError(ctx, err)
	}

	return ProcessDataHTTPResponse(ctx, http.StatusOK, getUserResponse{User: user})
}
