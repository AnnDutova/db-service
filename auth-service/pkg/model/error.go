package model

import (
	"errors"
	"fmt"
	"net/http"
)

type Error struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func NewError(message string, status int) *Error {
	return &Error{
		Message: message,
		Status:  status,
	}
}

func (e *Error) Error() string {
	return e.Message
}

func Status(err error) int {
	var e *Error
	if errors.As(err, &e) {
		return e.Status
	}
	return http.StatusInternalServerError
}

func UnauthorizedError(reason string) *Error {
	status := http.StatusUnauthorized
	return NewError(reason, status)
}

func BadRequestError(reason string) *Error {
	message := fmt.Sprintf("Bad request for a reason %s", reason)
	status := http.StatusBadRequest
	return NewError(message, status)
}

func InternalError() *Error {
	message := "Internal server error"
	status := http.StatusInternalServerError
	return NewError(message, status)
}

func NotFoundError(name string, value string) *Error {
	message := fmt.Sprintf("resource: %v with value: %v not found", name, value)
	status := http.StatusNotFound
	return NewError(message, status)
}

func ConflictError(name string, value string) *Error {
	message := fmt.Sprintf("resource: %v with value: %v already exists", name, value)
	status := http.StatusConflict
	return NewError(message, status)
}

func NewUnsupportedMediaType(message string) *Error {
	status := http.StatusUnsupportedMediaType
	return NewError(message, status)
}

func NewServiceUnavailable() *Error {
	status := http.StatusServiceUnavailable
	message := fmt.Sprintf("Service unavailable or timed out")
	return NewError(message, status)
}
