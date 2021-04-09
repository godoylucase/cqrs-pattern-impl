package internal

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrValueValidation  = errors.New("invalid value")
	ErrResourceNotFound = errors.New("resource not found")
)

type AppError struct {
	Cause error
	Type  error
}

func (ae *AppError) Error() string {
	if ae.Cause != nil {
		return fmt.Sprintf("an error type: %v with cause: %v", ae.Type, ae.Cause)
	}

	return fmt.Sprintf("an error type: %v", ae.Type)
}

func (ae *AppError) GetHTTPStatusFromType() int {
	switch ae.Type {
	case ErrResourceNotFound:
		return http.StatusNotFound
	case ErrValueValidation:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
