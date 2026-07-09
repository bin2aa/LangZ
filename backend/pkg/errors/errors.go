package errors

import "net/http"

type AppError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	HTTPStatus int    `json:"-"`
	Err        error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// Predefined error types
var (
	ErrNotFound     = &AppError{Code: "NOT_FOUND", Message: "Resource not found", HTTPStatus: http.StatusNotFound}
	ErrUnauthorized = &AppError{Code: "UNAUTHORIZED", Message: "Unauthorized", HTTPStatus: http.StatusUnauthorized}
	ErrForbidden    = &AppError{Code: "FORBIDDEN", Message: "Forbidden", HTTPStatus: http.StatusForbidden}
	ErrConflict     = &AppError{Code: "CONFLICT", Message: "Resource already exists", HTTPStatus: http.StatusConflict}
	ErrValidation   = &AppError{Code: "VALIDATION_ERROR", Message: "Validation failed", HTTPStatus: http.StatusBadRequest}
	ErrInternal     = &AppError{Code: "INTERNAL_ERROR", Message: "Internal server error", HTTPStatus: http.StatusInternalServerError}
	ErrInvalidInput = &AppError{Code: "INVALID_INPUT", Message: "Invalid input", HTTPStatus: http.StatusUnprocessableEntity}
)

func NewNotFound(msg string) *AppError {
	return &AppError{Code: "NOT_FOUND", Message: msg, HTTPStatus: http.StatusNotFound}
}

func NewConflict(msg string) *AppError {
	return &AppError{Code: "CONFLICT", Message: msg, HTTPStatus: http.StatusConflict}
}

func NewValidation(msg string) *AppError {
	return &AppError{Code: "VALIDATION_ERROR", Message: msg, HTTPStatus: http.StatusBadRequest}
}

func NewUnauthorized(msg string) *AppError {
	return &AppError{Code: "UNAUTHORIZED", Message: msg, HTTPStatus: http.StatusUnauthorized}
}

func NewForbidden(msg string) *AppError {
	return &AppError{Code: "FORBIDDEN", Message: msg, HTTPStatus: http.StatusForbidden}
}

func NewInternal(err error) *AppError {
	return &AppError{Code: "INTERNAL_ERROR", Message: "Internal server error", HTTPStatus: http.StatusInternalServerError, Err: err}
}
