package middleware

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/hoffax/prodapi/server/types"
)

func FiberCustomErrorHandler(c *fiber.Ctx, err error) error {
	if apiError, ok := err.(*types.ApiErrors); ok {
		return c.Status(fiber.StatusBadRequest).JSON(apiError.ToJSON())
	}

	//var validationErrors validator.ValidationErrors
	//if errors.As(err, &validationErrors) {
	//	errorMessages := make([]map[string]string, 0)
	//	for _, err := range validationErrors {
	//		errorMessages = append(errorMessages, map[string]string{
	//			"field":   err.Field(),
	//			"message": fmt.Sprintf("failed on %v %v %v validation", err.Kind(), err.Tag(), err.Param()),
	//		})
	//	}
	//
	//	return c.Status(fiber.StatusBadRequest).JSON(map[string]interface{}{
	//		"code":    "validation_error",
	//		"message": "validation errors",
	//		"errors":  errorMessages,
	//	})
	//}

	var e *fiber.Error
	if errors.As(err, &e) {
		fmt.Printf("fiber error: %v\n", e)
		return c.Status(e.Code).JSON(map[string]string{"message": e.Error()})
	}

	fmt.Printf("Not intercepted error: %v\n", err)
	// default error response
	return c.Status(fiber.StatusInternalServerError).JSON(map[string]string{"message": "internal server error"})
}
