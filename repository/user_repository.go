package repository

import (
	"github.com/hoffax/prodapi/dbmodel"
)

type UserRepository struct {
	// You can add DB here
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		// Initialize DB here
	}
}

func (r *UserRepository) GetByID(userID string) (*dbmodel.User, error) {
	// Fetch user from DB
	return nil, nil
}
