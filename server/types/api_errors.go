package types

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

// ApiErrors represents a custom API error
type ApiErrors struct {
	StatusCode int    `json:"-"`
	Code       string `json:"code"`
	Message    string `json:"message"`
}

// Implement the error interface
func (a ApiErrors) Error() string {
	return fmt.Sprintf("Code: %s, Message: %s", a.Code, a.Message)
}

// ToJSON returns a map that can be used for JSON responses
func (a ApiErrors) ToJSON() map[string]string {
	return map[string]string{
		"code":    a.Code,
		"message": a.Message,
	}
}

// NewInvalidParamsError creates a new ApiErrors for invalid parameters
func NewInvalidParamsError(message string) error {
	return &ApiErrors{
		StatusCode: fiber.StatusBadRequest,
		Code:       "invalid_params",
		Message:    message,
	}
}

// NewInvalidBodyError creates a new ApiErrors for an invalid body
func NewInvalidBodyError() error {
	return &ApiErrors{
		StatusCode: fiber.StatusBadRequest,
		Code:       "bad_formatted_body",
		Message:    "invalid json on request body",
	}
}

// NewNotFoundError creates a new ApiErrors for not found resources
func NewNotFoundError() error {
	return &ApiErrors{
		StatusCode: fiber.StatusNotFound,
		Code:       "resource_not_found",
		Message:    "resource not found",
	}
}
