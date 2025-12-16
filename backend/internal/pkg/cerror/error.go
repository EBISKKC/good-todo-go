package cerror

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Code       ErrorCode              `json:"code"`
	Message    string                 `json:"message"`
	Details    map[string]interface{} `json:"details,omitempty"`
	HTTPStatus int                    `json:"-"`
	Err        error                  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

type ErrorCode string

const (
	ErrCodeBadRequest          ErrorCode = "BAD_REQUEST"
	ErrCodeUnauthorized        ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden           ErrorCode = "FORBIDDEN"
	ErrCodeNotFound            ErrorCode = "NOT_FOUND"
	ErrCodeConflict            ErrorCode = "CONFLICT"
	ErrCodeInternalServerError ErrorCode = "INTERNAL_SERVER_ERROR"
	ErrCodeValidationError     ErrorCode = "VALIDATION_ERROR"
)

func NewBadRequest(message string, err error) *AppError {
	return &AppError{
		Code:       ErrCodeBadRequest,
		Message:    message,
		HTTPStatus: http.StatusBadRequest,
		Err:        err,
	}
}

func NewUnauthorized(message string, err error) *AppError {
	return &AppError{
		Code:       ErrCodeUnauthorized,
		Message:    message,
		HTTPStatus: http.StatusUnauthorized,
		Err:        err,
	}
}

func NewForbidden(message string, err error) *AppError {
	return &AppError{
		Code:       ErrCodeForbidden,
		Message:    message,
		HTTPStatus: http.StatusForbidden,
		Err:        err,
	}
}

func NewNotFound(message string, err error) *AppError {
	return &AppError{
		Code:       ErrCodeNotFound,
		Message:    message,
		HTTPStatus: http.StatusNotFound,
		Err:        err,
	}
}

func NewConflict(message string, err error) *AppError {
	return &AppError{
		Code:       ErrCodeConflict,
		Message:    message,
		HTTPStatus: http.StatusConflict,
		Err:        err,
	}
}

func NewInternalServerError(message string, err error) *AppError {
	return &AppError{
		Code:       ErrCodeInternalServerError,
		Message:    message,
		HTTPStatus: http.StatusInternalServerError,
		Err:        err,
	}
}

func NewValidationError(message string, details map[string]interface{}) *AppError {
	return &AppError{
		Code:       ErrCodeValidationError,
		Message:    message,
		HTTPStatus: http.StatusBadRequest,
		Details:    details,
	}
}
