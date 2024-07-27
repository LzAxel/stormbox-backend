package http

import (
	"chat-backend/internal/apperror"
	"strings"

	"chat-backend/internal/jwt"
	"github.com/labstack/echo/v4"
)

type JWTValidator interface {
	ValidateToken(token string) (jwt.Claims, error)
}

func (h *Handler) Authorized() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			h.logger.Debug("auth middleware", map[string]interface{}{"auth_header": authHeader})
			splitAuth := strings.Split(authHeader, " ")
			if len(splitAuth) != 2 {
				if c.QueryParam("access_token") == "" {
					return h.handleAppError(c, apperror.ErrInvalidAccessToken)
				}

				splitAuth = strings.Split("Bearer "+c.QueryParam("access_token"), " ")
				if (len(splitAuth) != 2) || (splitAuth[0] != "Bearer") {
					return h.handleAppError(c, apperror.ErrInvalidAccessToken)
				}
			}

			claims, err := h.jwtValidator.ValidateToken(splitAuth[1])
			if err != nil {
				return h.handleAppError(c, apperror.ErrInvalidAccessToken)
			}
			user, err := h.services.User.GetByID(c.Request().Context(), claims.Subject)
			if err != nil {
				return apperror.NewDatabaseError("Authorized.GetByID: failed to get user", err)
			}
			h.logger.Debug("auth middleware", map[string]interface{}{"user": user})

			c.Set("user", user)

			return next(c)
		}
	}
}

//func (h *Handler) RequireUserType(userTypes ...model.UserType) echo.MiddlewareFunc {
//	return func(next echo.HandlerFunc) echo.HandlerFunc {
//		return func(c echo.Context) error {
//			user, ok := c.Get("user").(model.User)
//			if !ok {
//				return h.newAppErrorResponse(c, errors.New("Handler.RequireUserType:failed to get user from context (forgot to use Authorize middleware)"))
//			}
//
//			if !slices.Contains(userTypes, user.Type) {
//				return h.newAuthErrorResponse(c, http.StatusForbidden, errors.New("access denied"))
//			}
//
//			if err := next(c); err != nil {
//				c.Error(err)
//			}
//
//			return nil
//		}
//	}
//}
