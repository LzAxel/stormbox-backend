package http

import (
	"chat-backend/internal/apperror"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

func appErrorTypeToCode(_type int) int {
	switch _type {
	case apperror.ErrorTypeDatabase:
		return http.StatusInternalServerError
	case apperror.ErrorTypeNotFound:
		return http.StatusNotFound
	case apperror.ErrorTypeConflict:
		return http.StatusConflict
	case apperror.ErrorTypeForbidden:
		return http.StatusForbidden
	case apperror.ErrorTypeUnauthorized:
		return http.StatusUnauthorized
	case apperror.ErrorTypeBadRequest:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func (h *Handler) handleAppError(ctx echo.Context, err error) error {
	if apperror.IsAppError(err) {
		var appErr *apperror.AppError
		errors.As(err, &appErr)
		status := appErrorTypeToCode(appErr.Type)

		h.logger.Debugf("request error: %s:%s", appErr.Message, appErr.Err)
		return ctx.JSON(status, ErrorHTTPResponse{
			Status:       status,
			ErrorMessage: appErr.Message,
		})
	}
	h.logger.Errorf("unexpected error: %s", err.Error())
	return ctx.JSON(http.StatusInternalServerError, ErrorHTTPResponse{
		Status:       http.StatusInternalServerError,
		ErrorMessage: "internal server error",
	})

}
