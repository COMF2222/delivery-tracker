package service

import (
	"delivery-tracker/internal/domain"
	"delivery-tracker/internal/repository"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Login(login, password string) (*domain.User, error) {
	user, err := s.userRepo.GetByLogin(login)
	if err != nil {
		return nil, fmt.Errorf("get by login: %w", err)
	}

	if !user.IsActive {
		return nil, ErrUserInactive
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidPassword
	}

	return user, nil
}
