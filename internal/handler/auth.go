package handler

import (
	"delivery-tracker/internal/dto"
	"delivery-tracker/internal/repository"
	"delivery-tracker/internal/service"
	"encoding/json"
	"errors"
	"net/http"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req dto.LoginUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "failed to decode request", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.authService.Login(req.Login, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidPassword) || errors.Is(err, repository.ErrUserNotFound) {
			http.Error(w, "invalid login or password", http.StatusUnauthorized)
			return
		}
		if errors.Is(err, service.ErrUserInactive) {
			http.Error(w, "user is inactive", http.StatusForbidden)
			return
		}
		http.Error(w, "failed to login user", http.StatusInternalServerError)
		return
	}

	resp := dto.LoginUserResponse{Login: user.Login}
	responseJSON, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(responseJSON)
	if err != nil {
		return
	}
}
