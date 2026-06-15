package service

import (
	"delivery-tracker/internal/domain"
	"delivery-tracker/internal/repository"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) CreateUser(login, password string, role domain.Role) (*domain.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return nil, fmt.Errorf("bcrypt password: %w", err)
	}

	user := domain.User{
		Login:        login,
		PasswordHash: string(hash),
		Role:         role,
	}

	err = s.userRepo.Create(&user)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return &user, nil
}
