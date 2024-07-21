package http

import (
	"chat-backend/internal/apperror"
	"chat-backend/internal/model"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (h *Handler) WithPagination() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			offsetParam := c.QueryParam("offset")
			if offsetParam == "" {
				offsetParam = "0"
			}
			offset, err := strconv.ParseUint(offsetParam, 10, 64)
			if err != nil {
				return h.handleAppError(c, apperror.ErrInvalidPaginationOffset)
			}

			limitParam := c.QueryParam("limit")
			if limitParam == "" {
				limitParam = strconv.Itoa(model.DefaultLimit)
			}
			limit, err := strconv.ParseUint(limitParam, 10, 64)
			if err != nil {
				return h.handleAppError(c, apperror.ErrInvalidPaginationLimit)
			}

			pagination, err := model.NewPagination(offset, limit)
			if err != nil {
				return h.handleAppError(c, err)
			}
			c.Set("pagination", pagination)

			h.logger.Debug("pagination", map[string]interface{}{
				"offset":     offset,
				"limit":      limit,
				"request_id": c.Response().Header().Get(echo.HeaderXRequestID),
			})
			return next(c)
		}
	}
}

func (h *Handler) getPaginationFromContext(ctx echo.Context) (model.Pagination, error) {
	pagination, ok := ctx.Get("pagination").(model.Pagination)
	if !ok {
		pagination = model.Pagination{
			Offset: 0,
			Limit:  model.DefaultLimit,
		}
	}
	h.logger.Debug("pagination", map[string]interface{}{
		"offset":     pagination.Offset,
		"limit":      pagination.Limit,
		"request_id": ctx.Response().Header().Get(echo.HeaderXRequestID),
	})
	return pagination, nil
}
