package service

import (
	"chat-backend/internal/apperror"
	"chat-backend/internal/jwt"
	"chat-backend/internal/model"
	"chat-backend/internal/repository"
	"chat-backend/pkg/clock"
	"chat-backend/pkg/hash"
	"context"
	"errors"
)

type AuthorizationService struct {
	jwt      *jwt.JWT
	userRepo repository.User
}

func NewAuthorizationService(jwt *jwt.JWT, userRepo repository.User) *AuthorizationService {
	return &AuthorizationService{
		jwt:      jwt,
		userRepo: userRepo,
	}
}

func (a *AuthorizationService) RefreshTokens(ctx context.Context, refreshToken string) (jwt.TokenPair, error) {
	claims, err := a.jwt.ValidateToken(refreshToken)
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrInvalidToken) || errors.Is(err, jwt.ErrInvalidClaims):
			return jwt.TokenPair{}, apperror.ErrInvalidRefreshToken
		case errors.Is(err, jwt.ErrTokenExpired):
			return jwt.TokenPair{}, apperror.ErrTokenExpired
		}
		return jwt.TokenPair{}, apperror.NewServiceError("RefreshTokens.ValidateToken", err)
	}

	tokenPair, err := a.jwt.GeneratePair(claims.Subject)
	if err != nil {
		return jwt.TokenPair{}, apperror.NewServiceError("RefreshTokens.GeneratePair", err)
	}

	return tokenPair, err
}

func (a *AuthorizationService) Login(ctx context.Context, input model.LoginInput) (jwt.TokenPair, error) {
	user, err := a.userRepo.GetByUsername(ctx, input.Username)
	if err != nil {
		if errors.Is(err, apperror.ErrUserNotFound) {
			return jwt.TokenPair{}, apperror.ErrInvalidLoginOrPassword
		}

		return jwt.TokenPair{}, apperror.NewDatabaseError("Login.GetByUsername", err)
	}

	if err := hash.Compare(user.PasswordHash, input.Password); err != nil {
		return jwt.TokenPair{}, apperror.ErrInvalidLoginOrPassword
	}

	tokenPair, err := a.jwt.GeneratePair(user.ID)
	if err != nil {
		return jwt.TokenPair{}, apperror.NewServiceError("Login.GeneratePair", err)
	}

	return tokenPair, err
}
func (a *AuthorizationService) Register(ctx context.Context, input model.CreateUserInput) error {
	passwordHash, err := hash.Hash(input.Password)
	if err != nil {
		return err
	}

	dto := model.CreateUserDTO{
		Username:     input.Username,
		PasswordHash: passwordHash,
		OnlineAt:     clock.Now().UTC(),
		CreatedAt:    clock.Now().UTC(),
		UpdatedAt:    clock.Now().UTC(),
	}

	if _, err := a.userRepo.Create(ctx, dto); err != nil {
		if apperror.IsAppError(err) {
			return err
		}

		return apperror.NewDatabaseError("Register.Create", err)
	}

	return nil
}
