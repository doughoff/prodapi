package routes

import (
	"bytes"
	"database/sql"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/hoffax/prodapi/postgres"
	"github.com/hoffax/prodapi/server/dto"
	"github.com/hoffax/prodapi/server/types"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type GetAllUsersQuery struct {
	StatusOptions []string `query:"status"`
	Search        string   `query:"search"`
	Limit         int32    `query:"limit"`
	Offset        int32    `query:"offset"`
}

func (r *RouteManager) getAllUsers(c *fiber.Ctx, tx *pgx.Tx) error {
	params := new(GetAllEntitiesQuery)
	if err := c.QueryParser(params); err != nil {
		return types.NewInvalidParamsError("invalid query params")
	}

	statusOptions := make([]postgres.Status, len(params.StatusOptions))
	for i, status := range params.StatusOptions {
		statusOptions[i] = postgres.Status(status)
		if !statusOptions[i].Valid() {
			return types.NewInvalidParamsError("invalid status option")
		}
	}

	users, err := r.db.GetUsers(c.Context(), *tx, &postgres.GetUsersParams{
		PageOffset:    params.Offset,
		PageLimit:     params.Limit,
		Search:        pgtype.Text(sql.NullString{String: params.Search, Valid: true}),
		StatusOptions: statusOptions,
	})
	if err != nil {
		return err
	}

	var totalCount int64
	if len(users) > 0 {
		totalCount = users[0].FullCount
	}
	resultRows := make([]*dto.UserDTO, len(users))
	for i := range resultRows {
		resultRows[i] = &dto.UserDTO{
			ID:        users[i].ID,
			Status:    users[i].Status,
			Name:      users[i].Name,
			Email:     users[i].Email,
			Roles:     users[i].Roles,
			CreatedAt: users[i].CreatedAt.Time,
			UpdatedAt: users[i].UpdatedAt.Time,
		}
	}

	return c.JSON(struct {
		TotalCount int64          `json:"totalCount"`
		Items      []*dto.UserDTO `json:"items"`
	}{
		TotalCount: totalCount,
		Items:      resultRows,
	})
}

type CreateUserBody struct {
	Name     string   `json:"name" validate:"required,gte=3,lte=255"`
	Email    string   `json:"email" validate:"required,email"`
	Roles    []string `json:"roles"`
	Password string   `json:"password" validate:"required,gte=8,lte=255"`
}

func (r *RouteManager) createUser(c *fiber.Ctx, tx *pgx.Tx) error {
	body := new(CreateUserBody)
	if err := c.BodyParser(body); err != nil {
		return types.NewInvalidBodyError()
	}
	if err := r.validate.Struct(body); err != nil {
		return err
	}

	_, err := r.db.GetUserByEmail(c.Context(), *tx, body.Email)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return err
		}
	} else {
		return types.NewInvalidParamsError("user with this email already exists")
	}

	user, err := r.db.CreateUser(c.Context(), *tx, &postgres.CreateUserParams{
		Name:     body.Name,
		Email:    body.Email,
		Roles:    body.Roles,
		Password: body.Password,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(dto.ToUserDTO(user))
}

type UpdateUserBody struct {
	Status postgres.Status `json:"status" validate:"required"`
	Name   string          `json:"name" validate:"omitempty,gte=3,lte=255"`
	Email  string          `json:"email" validate:"omitempty,email"`
	Roles  []string        `json:"roles"`
}

func (r *RouteManager) updateUserByID(c *fiber.Ctx, tx *pgx.Tx) error {
	idParam := c.Params("id")
	userID := pgtype.UUID{}
	if err := userID.Scan(idParam); err != nil {
		return types.NewInvalidParamsError("invalid uuid on id url param")
	}

	body := new(UpdateUserBody)
	if err := c.BodyParser(body); err != nil {
		return types.NewInvalidBodyError()
	}
	if err := r.validate.Struct(body); err != nil {
		return err
	}

	if !body.Status.Valid() {
		return types.NewInvalidParamsError("invalid value for status")
	}

	user, err := r.db.GetUserByEmail(c.Context(), *tx, body.Email)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return err
		}
	} else {
		if bytes.Equal(user.ID.Bytes[:], userID.Bytes[:]) {
			return types.NewInvalidParamsError("user with this email already exists")
		}
	}

	user, err = r.db.UpdateUserByID(c.Context(), *tx, &postgres.UpdateUserByIDParams{
		ID:     userID,
		Status: body.Status,
		Name:   body.Name,
		Email:  body.Email,
		Roles:  body.Roles,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dto.ToUserDTO(user))
}

type UpdateUserPasswordBody struct {
	Password string `json:"password" validate:"required,gte=8,lte=255"`
}

func (r *RouteManager) updateUserPasswordByID(c *fiber.Ctx, tx *pgx.Tx) error {
	idParam := c.Params("id")
	userID := pgtype.UUID{}
	if err := userID.Scan(idParam); err != nil {
		return types.NewInvalidParamsError("invalid uuid on id url param")
	}

	body := new(UpdateUserPasswordBody)
	if err := c.BodyParser(body); err != nil {
		return types.NewInvalidBodyError()
	}
	if err := r.validate.Struct(body); err != nil {
		return err
	}

	user, err := r.db.UpdateUserByID(c.Context(), *tx, &postgres.UpdateUserByIDParams{
		ID:       userID,
		Password: body.Password,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dto.ToUserDTO(user))
}

func (r *RouteManager) getUserByID(c *fiber.Ctx, tx *pgx.Tx) error {
	idParam := c.Params("id")
	userID := pgtype.UUID{}
	if err := userID.Scan(idParam); err != nil {
		return types.NewInvalidParamsError("invalid uuid on id url param")
	}

	user, err := r.db.GetUserByID(c.Context(), *tx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return types.NewNotFoundError()
		}
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dto.ToUserDTO(user))
}

func (r *RouteManager) getUserByEmail(c *fiber.Ctx, tx *pgx.Tx) error {
	emailParam := c.Params("email")
	user, err := r.db.GetUserByEmail(c.Context(), *tx, emailParam)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return types.NewNotFoundError()
		}
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dto.ToUserDTO(user))
}

func (r *RouteManager) RegisterUserRoutes() {
	g := r.app.Group("/users")

	g.Get("/", r.dbWrapper.WithTransaction(r.getAllUsers))
	g.Post("/", r.dbWrapper.WithTransaction(r.createUser))
	g.Put("/:id", r.dbWrapper.WithTransaction(r.updateUserByID))
	g.Put("/:id/reset_password", r.dbWrapper.WithTransaction(r.updateUserPasswordByID))
	g.Get("/:id", r.dbWrapper.WithTransaction(r.getUserByID))
	r.app.Get("/check_email/:email", r.dbWrapper.WithTransaction(r.getUserByEmail))

}
