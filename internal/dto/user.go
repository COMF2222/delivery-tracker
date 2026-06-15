package dto

import (
	"delivery-tracker/internal/domain"
	"fmt"
)

type CreateUserRequest struct {
	Login    string      `json:"login"`
	Password string      `json:"password"`
	Role     domain.Role `json:"role"`
}

type CreateUserResponse struct {
	ID int `json:"id"`
}

func (r CreateUserRequest) Validate() error {
	if r.Login == "" {
		return fmt.Errorf("login cannot be empty")
	}

	if r.Password == "" {
		return fmt.Errorf("password cannot be empty")
	}

	if r.Role != domain.RoleAdmin &&
		r.Role != domain.RoleManager {

		return fmt.Errorf("invalid role")
	}

	return nil
}
