package routes

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/hoffax/prodapi/server/types"
	"github.com/jackc/pgx/v5"
)

type LoginPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (r *RouteManager) login(c *fiber.Ctx, tx *pgx.Tx) error {

	body := new(LoginPayload)
	if err := c.BodyParser(body); err != nil {
		return types.NewInvalidBodyError()
	}
	if err := r.validate.Struct(body); err != nil {
		return err
	}

	user, err := r.db.GetUserByEmail(c.Context(), *tx, body.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusBadRequest).JSON(map[string]string{"message": "invalid email or password"})
		}
		return err
	}

	if user.Password != body.Password {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]string{"message": "invalid email or password"})
	}

	sess, err := r.sessionStore.Get(c)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	userIDStr, err := user.ID.MarshalJSON()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	sess.Set("userID", userIDStr)

	if err := sess.Save(); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).Send([]byte("ok"))
}

func (r *RouteManager) logout(c *fiber.Ctx, _ *pgx.Tx) error {
	sess, err := r.sessionStore.Get(c)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Destroy session
	if err := sess.Destroy(); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(map[string]string{"message": "logged out"})
}

func (r *RouteManager) RegisterAuthRoutes() {
	r.app.Post("/auth/login", r.dbWrapper.WithTransaction(r.login))
	r.app.Post("/auth/logout", r.dbWrapper.WithTransaction(r.logout))
}
