package httperror

import (
	"errors"
	"net/http"
)

type HTTPError struct {
	error
	Code int
}

func New(code int, message string) *HTTPError {
	return &HTTPError{
		error: errors.New(message),
		Code:  code,
	}
}

func NotFound(message string) *HTTPError {
	return New(http.StatusNotFound, message)
}
