package apperror

import "strconv"

var (
	ErrInvalidPaginationLimit  = New(ErrorTypeBadRequest, "invalid pagination limit", nil)
	ErrInvalidPaginationOffset = New(ErrorTypeBadRequest, "invalid pagination offset", nil)

	// User
	ErrUserNotFound          = New(ErrorTypeNotFound, "user not found", nil)
	ErrUsernameAlreadyExists = New(ErrorTypeConflict, "username already exists", nil)

	// Auth
	ErrInvalidAccessToken     = New(ErrorTypeUnauthorized, "invalid access token", nil)
	ErrInvalidRefreshToken    = New(ErrorTypeBadRequest, "invalid refresh token", nil)
	ErrTokenExpired           = New(ErrorTypeUnauthorized, "token expired", nil)
	ErrInvalidLoginOrPassword = New(ErrorTypeUnauthorized, "invalid login or password", nil)
)

func GetErrMaxPaginationLimit(maxLimit uint64) error {
	return New(ErrorTypeBadRequest, "max pagination limit is "+strconv.FormatUint(maxLimit, 10), nil)
}
