package http

import (
	"chat-backend/internal/apperror"
	"net/http"

	"chat-backend/internal/model"

	"github.com/labstack/echo/v4"
)

type signUpRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type signUpResponse struct {
	User  model.ViewUser `json:"user"`
	Token tokensResponse `json:"token"`
}

func (h *Handler) signUp(ctx echo.Context) error {
	var req signUpRequest

	if err := ctx.Bind(&req); err != nil {
		return apperror.New(apperror.ErrorTypeBadRequest, "invalid input", err)
	}

	user, tokenPair, err := h.services.Authorization.Register(ctx.Request().Context(), model.CreateUserInput{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return h.handleAppError(ctx, err)
	}

	return ProcessDataHTTPResponse(ctx, http.StatusOK, signUpResponse{
		User: user,
		Token: tokensResponse{
			AccessToken:  tokenPair.AccessToken,
			RefreshToken: tokenPair.RefreshToken,
		},
	})
}

type signInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type tokensResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type signInResponse struct {
	User  model.ViewUser `json:"user"`
	Token tokensResponse `json:"token"`
}

func (h *Handler) signIn(ctx echo.Context) error {
	var req signInRequest

	if err := ctx.Bind(&req); err != nil {
		return apperror.New(apperror.ErrorTypeBadRequest, "invalid input", err)
	}

	user, tokenPair, err := h.services.Authorization.Login(ctx.Request().Context(), model.LoginInput{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return h.handleAppError(ctx, err)
	}

	return ProcessDataHTTPResponse(ctx, http.StatusOK, signInResponse{
		User: user,
		Token: tokensResponse{
			AccessToken:  tokenPair.AccessToken,
			RefreshToken: tokenPair.RefreshToken,
		},
	})
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *Handler) refreshTokens(ctx echo.Context) error {
	var req refreshRequest

	if err := ctx.Bind(&req); err != nil {
		return apperror.New(apperror.ErrorTypeBadRequest, "invalid input", err)
	}

	tokenPair, err := h.services.Authorization.RefreshTokens(ctx.Request().Context(), req.RefreshToken)
	if err != nil {
		return h.handleAppError(ctx, err)
	}
	return ProcessDataHTTPResponse(ctx, http.StatusOK, tokensResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	})
}
