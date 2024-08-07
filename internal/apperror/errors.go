package apperror

import (
	"fmt"
	"strconv"
)

var (
	ErrInvalidPaginationLimit  = New(ErrorTypeBadRequest, "invalid pagination limit", nil)
	ErrInvalidPaginationOffset = New(ErrorTypeBadRequest, "invalid pagination offset", nil)

	// User
	ErrUserNotFound          = New(ErrorTypeNotFound, "user not found", nil)
	ErrUsernameAlreadyExists = New(ErrorTypeConflict, "username already exists", nil)

	// Auth
	ErrInvalidAccessToken     = New(ErrorTypeUnauthorized, "invalid access token", nil)
	ErrInvalidRefreshToken    = New(ErrorTypeBadRequest, "invalid refresh token", nil)
	ErrAccessTokenExpired     = New(ErrorTypeUnauthorized, "access token expired", nil)
	ErrRefreshTokenExpired    = New(ErrorTypeUnauthorized, "refresh token expired", nil)
	ErrInvalidLoginOrPassword = New(ErrorTypeUnauthorized, "invalid login or password", nil)

	// Friendship
	ErrFriendshipAlreadyExists = New(ErrorTypeConflict, "friendship already exists", nil)
	ErrCannotFriendSelf        = New(ErrorTypeBadRequest, "cannot create friendship with yourself", nil)

	// Message
	ErrCanSendMessagesOnlyFriends = New(ErrorTypeForbidden, "can send messages only to friends", nil)

	// Handler
	ErrInvalidRequestBody = New(ErrorTypeBadRequest, "invalid request body", nil)
)

func GetErrInvalidParam(param string) error {
	return New(ErrorTypeBadRequest, fmt.Sprintf("invalid '%s' value", param), nil)
}

func GetErrMaxPaginationLimit(maxLimit uint64) error {
	return New(ErrorTypeBadRequest, "max pagination limit is "+strconv.FormatUint(maxLimit, 10), nil)
}
