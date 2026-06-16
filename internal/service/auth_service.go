package service

import (
	"delivery-tracker/internal/auth"
	"delivery-tracker/internal/repository"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type AuthService struct {
	userRepo *repository.UserRepository
	secret   string
	ttl      time.Duration
}

func NewAuthService(userRepo *repository.UserRepository, secret string, ttl time.Duration) *AuthService {
	return &AuthService{userRepo: userRepo, secret: secret, ttl: ttl}
}

func (s *AuthService) Login(login, password string) (string, error) {
	user, err := s.userRepo.GetByLogin(login)
	if err != nil {
		return "", fmt.Errorf("get by login: %w", err)
	}

	if !user.IsActive {
		return "", ErrUserInactive
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", ErrInvalidPassword
	}

	token, err := auth.GenerateToken(user.ID, user.Login, user.Role, s.secret, s.ttl)
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}

	return token, nil
}
