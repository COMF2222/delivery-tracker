package domain

import "time"

type User struct {
	ID           int       `db:"id"`
	Login        string    `db:"login"`
	PasswordHash string    `db:"password_hash"`
	Role         Role      `db:"role"`
	IsActive     bool      `db:"is_active"`
	CreatedAt    time.Time `db:"created_at"`
}

type Role string

const (
	RoleAdmin   Role = "admin"
	RoleManager Role = "manager"
)
