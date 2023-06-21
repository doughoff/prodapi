package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func AuthMiddleware(store *session.Store) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		if c.Path() == "/auth/login" {
			return c.Next()
		}
		sess, err := store.Get(c)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		userID := sess.Get("userID")
		if userID == nil {
			return fiber.NewError(fiber.StatusUnauthorized, "unauthenticated")
		}

		c.Locals("userId", userID)

		return c.Next()
	}
}
