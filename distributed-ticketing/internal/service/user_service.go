package service

import (
	"distributed-ticketing/internal/model"
	"distributed-ticketing/internal/repository"
)

type UserService struct {
	userRepository *repository.UserRepository
}

func NewUserService(userRepository *repository.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) CreateUser(email string) (*model.User, error) {
	return s.userRepository.CreateUser(email)
}
func (s *UserService) GetUserByEmail(email string) (*model.User, error) {
	return s.userRepository.GetUserByEmail(email)
}
