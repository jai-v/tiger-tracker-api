package apiError

import (
	"fmt"
	"net/http"
)

type ErrorCode string

func (e ErrorCode) String() string {
	return string(e)
}

const (
	BadRequestErrorCode                       ErrorCode = "ERR_BAD_REQUEST"
	InternalServerErrorCode                             = "ERR_INTERNAL_SERVER_ERROR"
	UsernameOrEmailAlreadyRegisteredErrorCode           = "ERR_USERNAME_OR_EMAIL_ALREADY_REGISTERED"
	UnauthorizedErrorCode                               = "ERR_UNAUTHORIZED"
)

type APIError struct {
	ErrorCode      ErrorCode `json:"errorCode"`
	ErrorMessage   string    `json:"errorMessage"`
	HttpStatusCode int       `json:"-"`
}

func NewApiError(errorCode ErrorCode, errorMessage string, httpStatusCode int) *APIError {
	return &APIError{ErrorCode: errorCode, ErrorMessage: errorMessage, HttpStatusCode: httpStatusCode}
}

func (e *APIError) Error() string {
	return fmt.Sprintf("errorCode: %s , errorMessage: %s", e.ErrorCode, e.ErrorMessage)
}

var (
	BadRequestError                  = NewApiError(BadRequestErrorCode, "Invalid Request", http.StatusBadRequest)
	InternalServerError              = NewApiError(InternalServerErrorCode, "Internal Server Error", http.StatusInternalServerError)
	UsernameOrEmailAlreadyRegistered = NewApiError(UsernameOrEmailAlreadyRegisteredErrorCode, "Username or Email already registered", http.StatusUnprocessableEntity)
	UnauthorizedError                = NewApiError(UnauthorizedErrorCode, "Unauthorized", http.StatusUnauthorized)
)

func NewBadRequestErrorWithMessage(errorMessage string) *APIError {
	return NewApiError(BadRequestErrorCode, errorMessage, http.StatusBadRequest)
}
