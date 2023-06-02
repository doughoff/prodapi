package user

import (
	"context"
	"github.com/hoffax/prodapi/dbtypes"
	uuid "github.com/jackc/pgx-gofrs-uuid"
)

type Repository interface {
	NewUser(ctx context.Context, user *dbtypes.User) (*dbtypes.User, error)
	GetByID(ctx context.Context, userID uuid.UUID) (*dbtypes.User, error)
	All(ctx context.Context) ([]*dbtypes.User, error)
}
