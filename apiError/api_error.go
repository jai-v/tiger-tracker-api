package apiError

import "fmt"

type ErrorCode string

func (e ErrorCode) String() string {
	return string(e)
}

const (
	BadRequestErrorCode     ErrorCode = "ERR_BAD_REQUEST"
	InternalServerErrorCode ErrorCode = "ERR_INTERNAL_SERVER_ERROR"
)

type APIError struct {
	ErrorCode    ErrorCode `json:"errorCode"`
	ErrorMessage string    `json:"errorMessage"`
}

func NewApiError(errorCode ErrorCode, errorMessage string) *APIError {
	return &APIError{ErrorCode: errorCode, ErrorMessage: errorMessage}
}

func (e *APIError) Error() string {
	return fmt.Sprintf("errorCode: %s , errorMessage: %s", e.ErrorCode, e.ErrorMessage)
}

var (
	BadRequestError     = NewApiError(BadRequestErrorCode, "Invalid Request")
	InternalServerError = NewApiError(InternalServerErrorCode, "Internal Server Error")
)
