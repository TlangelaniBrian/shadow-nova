package errors

import (
	"errors"
	"fmt"
)

// Sentinel errors for common cases
var (
	ErrNotFound       = errors.New("resource not found")
	ErrUnauthorized   = errors.New("unauthorized access")
	ErrForbidden      = errors.New("forbidden")
	ErrInvalidInput   = errors.New("invalid input")
	ErrDuplicateEntry = errors.New("duplicate entry")
	ErrDatabaseError  = errors.New("database error")
)

// AppError is a custom error type with context
type AppError struct {
	Err      error
	Message  string
	Code     string
	HTTPCode int
}

func (e *AppError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Err.Error()
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// Error constructors
func NotFound(message string) *AppError {
	return &AppError{
		Err:      ErrNotFound,
		Message:  message,
		Code:     "NOT_FOUND",
		HTTPCode: 404,
	}
}

func Unauthorized(message string) *AppError {
	return &AppError{
		Err:      ErrUnauthorized,
		Message:  message,
		Code:     "UNAUTHORIZED",
		HTTPCode: 401,
	}
}

func Forbidden(message string) *AppError {
	return &AppError{
		Err:      ErrForbidden,
		Message:  message,
		Code:     "FORBIDDEN",
		HTTPCode: 403,
	}
}

func InvalidInput(message string) *AppError {
	return &AppError{
		Err:      ErrInvalidInput,
		Message:  message,
		Code:     "INVALID_INPUT",
		HTTPCode: 400,
	}
}

func DuplicateEntry(message string) *AppError {
	return &AppError{
		Err:      ErrDuplicateEntry,
		Message:  message,
		Code:     "DUPLICATE_ENTRY",
		HTTPCode: 409,
	}
}

func DatabaseError(err error, message string) *AppError {
	return &AppError{
		Err:      fmt.Errorf("%w: %v", ErrDatabaseError, err),
		Message:  message,
		Code:     "DATABASE_ERROR",
		HTTPCode: 500,
	}
}

// Check if error is specific type
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

func IsUnauthorized(err error) bool {
	return errors.Is(err, ErrUnauthorized)
}

func IsForbidden(err error) bool {
	return errors.Is(err, ErrForbidden)
}

func IsInvalidInput(err error) bool {
	return errors.Is(err, ErrInvalidInput)
}

func IsDuplicateEntry(err error) bool {
	return errors.Is(err, ErrDuplicateEntry)
}

func IsDatabaseError(err error) bool {
	return errors.Is(err, ErrDatabaseError)
}
