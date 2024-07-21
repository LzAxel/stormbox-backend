package apperror

import (
	"errors"
	"fmt"
)

const (
	ErrorTypeDatabase = iota
	ErrorTypeNotFound
	ErrorTypeConflict
	ErrorTypeForbidden
	ErrorTypeUnauthorized
	ErrorTypeBadRequest
)

type AppError struct {
	Type    int
	Message string
	Err     error
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(_type int, message string, err error) *AppError {
	return &AppError{
		Type:    _type,
		Message: message,
		Err:     err,
	}
}

func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}

func NewDatabaseError(message string, err error) error {
	return fmt.Errorf("database error: %s: %w", message, err)
}

func NewServiceError(message string, err error) error {
	return fmt.Errorf("service error: %s: %w", message, err)
}
