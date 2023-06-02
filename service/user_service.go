package service

import (
	"github.com/hoffax/prodapi/dbmodel"
	"github.com/hoffax/prodapi/repository"
)

type UserService struct {
	UserRepository *repository.UserRepository
}

func NewUserService(userRepository *repository.UserRepository) *UserService {
	return &UserService{
		UserRepository: userRepository,
	}
}

func (s *UserService) GetByID(userID string) (*dbmodel.User, error) {
	return s.UserRepository.GetByID(userID)
}
