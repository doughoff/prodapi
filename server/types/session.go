package types

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type SessionData struct {
	UserId pgtype.UUID
	Roles  []string
}
