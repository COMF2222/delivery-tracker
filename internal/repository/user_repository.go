package repository

import (
	"database/sql"
	"delivery-tracker/internal/domain"
	"errors"
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

func (r *UserRepository) GetByLogin(login string) (*domain.User, error) {
	var user domain.User
	query := `SELECT id, login, password_hash, role, is_active, created_at FROM users WHERE login = $1`

	if err := r.db.QueryRow(query, login).Scan(
		&user.ID,
		&user.Login,
		&user.PasswordHash,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by login: %w", err)
	}

	return &user, nil
}
