package repository

import (
	"delivery-tracker/internal/domain"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *domain.User) error {
	query := `INSERT INTO users(login, password_hash, role) VALUES ($1, $2, $3) RETURNING id`

	err := r.db.QueryRow(query, user.Login, user.PasswordHash, user.Role).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}
