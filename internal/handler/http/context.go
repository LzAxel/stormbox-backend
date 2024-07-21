package http

import (
	"chat-backend/internal/apperror"
	"chat-backend/internal/model"
	"github.com/labstack/echo/v4"
)

func getUserFromContext(ctx echo.Context) (model.ViewUser, error) {
	user, ok := ctx.Get("user").(model.ViewUser)
	if !ok {
		return model.ViewUser{}, apperror.NewServiceError("failed to get user from context", nil)
	}
	return user, nil
}
