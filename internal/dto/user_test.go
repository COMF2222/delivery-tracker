package dto

import (
	"delivery-tracker/internal/domain"
	"testing"
)

func TestCreateUserRequest_Validate(t *testing.T) {
	tests := []struct {
		name        string
		login       string
		password    string
		role        domain.Role
		expectError bool
	}{
		{name: "empty login", login: "", password: "pass", role: domain.RoleAdmin, expectError: true},
		{name: "empty password", login: "admin", password: "", role: domain.RoleAdmin, expectError: true},
		{name: "invalid role", login: "man", password: "pass", role: "", expectError: true},
		{name: "valid admin", login: "admin", password: "pass", role: domain.RoleAdmin, expectError: false},
		{name: "valid manager", login: "manager", password: "pass", role: domain.RoleManager, expectError: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := CreateUserRequest{
				Login:    tt.login,
				Password: tt.password,
				Role:     tt.role,
			}
			err := req.Validate()

			if tt.expectError {
				if err == nil {
					t.Fatalf("expected validation error for %#v, got nil", req)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected validation error: %+v: %v", req, err)
			}
		})
	}
}

func TestLoginUserRequest_Validate(t *testing.T) {
	tests := []struct {
		name        string
		login       string
		password    string
		expectError bool
	}{
		{name: "empty login", login: "", password: "pass", expectError: true},
		{name: "empty password", login: "login", password: "", expectError: true},
		{name: "valid login", login: "login", password: "pass", expectError: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := LoginUserRequest{
				Login:    tt.login,
				Password: tt.password,
			}
			err := req.Validate()

			if tt.expectError {
				if err == nil {
					t.Fatalf("expected validation error for %#v, got nil", req)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected validation error: %+v: %v", req, err)
			}
		})
	}
}
