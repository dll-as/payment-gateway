package errors

import (
	"github.com/gofiber/fiber/v3"
	"github.com/rezatg/payment-gateway/internal/models"
)

// ErrorCode represents standard error codes
type ErrorCode string

const (
	ErrBadRequest      ErrorCode = "BAD_REQUEST"
	ErrUnauthorized    ErrorCode = "UNAUTHORIZED"
	ErrNotFound        ErrorCode = "NOT_FOUND"
	ErrInternalServer  ErrorCode = "INTERNAL_SERVER"
	ErrConflict        ErrorCode = "CONFLICT"
	ErrTooManyRequests ErrorCode = "TOO_MANY_REQUESTS"
)

// APIError represents a structured API error
type APIError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"error"`
	Details string    `json:"details,omitempty"`
}

// Error implements the error interface
func (e *APIError) Error() string {
	return e.Message
}

// NewBadRequestError creates a new bad request error
func NewBadRequestError(message string, details ...string) *APIError {
	return newError(ErrBadRequest, message, details...)
}

// NewUnauthorizedError creates a new unauthorized error
func NewUnauthorizedError(message string, details ...string) *APIError {
	return newError(ErrUnauthorized, message, details...)
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(message string, details ...string) *APIError {
	return newError(ErrNotFound, message, details...)
}

// NewInternalServerError creates a new internal server error
func NewInternalServerError(message string, details ...string) *APIError {
	return newError(ErrInternalServer, message, details...)
}

// NewConflictError creates a new conflict error
func NewConflictError(message string, details ...string) *APIError {
	return newError(ErrConflict, message, details...)
}

// NewTooManyRequestsError creates a new too many requests error
func NewTooManyRequestsError(message string, details ...string) *APIError {
	return newError(ErrTooManyRequests, message, details...)
}

// newError creates a new APIError instance
func newError(code ErrorCode, message string, details ...string) *APIError {
	var detail string
	if len(details) > 0 {
		detail = details[0]
	}
	return &APIError{
		Code:    code,
		Message: message,
		Details: detail,
	}
}

// HandleAPIError converts an error to a Fiber response
func HandleAPIError(c fiber.Ctx, err error) error {
	if apiErr, ok := err.(*APIError); ok {
		return c.Status(getHTTPStatus(apiErr.Code)).JSON(models.ErrorResponse{
			Error:   apiErr.Message,
			Code:    getHTTPStatus(apiErr.Code),
			Details: apiErr.Details,
		})
	}
	return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
		Error: "Internal server error",
		Code:  fiber.StatusInternalServerError,
	})
}

// getHTTPStatus maps error codes to HTTP status codes
func getHTTPStatus(code ErrorCode) int {
	switch code {
	case ErrBadRequest:
		return fiber.StatusBadRequest
	case ErrUnauthorized:
		return fiber.StatusUnauthorized
	case ErrNotFound:
		return fiber.StatusNotFound
	case ErrConflict:
		return fiber.StatusConflict
	case ErrTooManyRequests:
		return fiber.StatusTooManyRequests
	default:
		return fiber.StatusInternalServerError
	}
}
