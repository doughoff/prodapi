package session

import "github.com/hoffax/prodapi/dbtypes"

type Repository interface {
	// You can add DB here
	GetByID(userID string) (*dbtypes.User, error)
}
