package service

import (
	"delivery-tracker/internal/domain"
	"delivery-tracker/internal/repository"
	"fmt"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	"strconv"
)

type UserService struct {
	userRepo  *repository.UserRepository
	auditRepo *repository.AuditRepository

	txManager *repository.TransactionManager
}

func NewUserService(userRepo *repository.UserRepository,
	auditRepo *repository.AuditRepository,
	txManager *repository.TransactionManager) *UserService {
	return &UserService{userRepo: userRepo, auditRepo: auditRepo, txManager: txManager}
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

func (s *UserService) Deactivate(userID, changedBy int) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("get user id: %w", err)
	}

	if !user.IsActive {
		return ErrUserAlreadyInactive
	}
	
	err = s.txManager.Do(func(tx *sqlx.Tx) error {
		oldValue := user.IsActive

		if err := s.userRepo.DeactivateTx(tx, userID); err != nil {
			return fmt.Errorf("failed to deactivate user: %w", err)
		}

		auditLog := domain.AuditLog{
			UserID:     changedBy,
			Action:     domain.ActionDeactivateUser,
			OldValue:   strconv.FormatBool(oldValue),
			NewValue:   strconv.FormatBool(false),
			EntityType: domain.EntityTypeUser,
			EntityID:   user.ID,
		}

		if err = s.auditRepo.CreateTx(tx, &auditLog); err != nil {
			return fmt.Errorf("failed to create audit log: %w", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("deactivate user transaction: %w", err)
	}

	return nil
}
