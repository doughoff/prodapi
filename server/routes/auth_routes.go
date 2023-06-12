package routes

import (
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/hoffax/prodapi/server/types"
	"github.com/jackc/pgx/v5"
	"time"
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

	// generate new uuid for session key
	sessionId, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	sessionData := &types.SessionData{
		UserId: user.ID,
		Roles:  user.Roles,
	}

	sessionDataBytes, err := json.Marshal(sessionData)
	if err != nil {
		return err
	}

	err = r.sessionStore.Set(sessionId.String(), sessionDataBytes, 72*time.Hour)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(map[string]string{"session_id": sessionId.String()})
}

func (r *RouteManager) logout(c *fiber.Ctx, tx *pgx.Tx) error {
	headers := c.GetReqHeaders()
	sessionID := headers["X-Session"]

	err := r.sessionStore.Delete(sessionID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(map[string]string{"message": "logged out"})
}

func (r *RouteManager) RegisterAuthRoutes() {
	r.app.Post("/auth/login", r.dbWrapper.WithTransaction(r.login))
	r.app.Post("/auth/logout", r.dbWrapper.WithTransaction(r.logout))
}
