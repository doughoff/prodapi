package dbtypes

import (
	uuid "github.com/jackc/pgx-gofrs-uuid"
	"time"
)

type User struct {
	ID        *uuid.UUID
	Status    string
	Roles     []string
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
