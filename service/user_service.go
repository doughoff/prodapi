package service

import (
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/hoffax/prodapi/dbtypes"
	"github.com/hoffax/prodapi/repository/user"
)

type UserService struct {
	UserRepository *user.PgRepository
}

func NewUserService(userRepository *user.PgRepository) *UserService {
	return &UserService{
		UserRepository: userRepository,
	}
}

func (s *UserService) GetByID(userID string) (*dbtypes.User, error) {
	return nil, nil
}
