package error

import (
	"fmt"
	"net/http"
)

type AppError interface {
	WithCode(code int) AppError
	WithReason(reason string) AppError

	Error() string
	ToHTTP() (statusCode int, message string)
}

var uniqueErrors = make(map[string]string)

func New(key, message string) AppError {
	uniqueErrors[key] = message

	return &appError{
		code:    http.StatusInternalServerError,
		message: message,
		descriptor: &descriptor{
			domain: key,
			reason: message,
		},
	}
}

type (
	appError struct {
		code       int
		message    string
		descriptor *descriptor
	}

	descriptor struct {
		domain string
		reason string
	}
)

func (e *appError) WithCode(code int) AppError {
	e.code = code
	return e
}

func (e *appError) WithReason(reason string) AppError {
	e.descriptor.reason = reason
	return e
}

func (e *appError) Error() string {
	return fmt.Sprintf("%s: %s", e.descriptor.domain, e.descriptor.reason)
}

func (e *appError) ToHTTP() (statusCode int, message string) {
	return e.code, e.message
}
