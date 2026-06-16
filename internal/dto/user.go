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

type LoginUserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginUserResponse struct {
	Token string `json:"token"`
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

func (r LoginUserRequest) Validate() error {
	if r.Login == "" {
		return fmt.Errorf("login cannot be empty")
	}

	if r.Password == "" {
		return fmt.Errorf("password cannot be empty")
	}

	return nil
}
