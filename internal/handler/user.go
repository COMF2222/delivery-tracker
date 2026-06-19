package handler

import (
	"delivery-tracker/internal/contextkeys"
	"delivery-tracker/internal/dto"
	"delivery-tracker/internal/repository"
	"delivery-tracker/internal/service"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req dto.CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "failed to decode request", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.userService.CreateUser(req.Login, req.Password, req.Role)
	if err != nil {
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	resp := dto.CreateUserResponse{ID: user.ID}

	responseJSON, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(responseJSON)
	if err != nil {
		return
	}
}

func (h *UserHandler) Deactivate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PATCH" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id cannot be empty", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "failed convert id to int", http.StatusBadRequest)
		return
	}
	if userID <= 0 {
		http.Error(w, "id must be positive", http.StatusBadRequest)
		return
	}

	changedBy, ok := r.Context().Value(contextkeys.UserID).(int)
	if !ok {
		http.Error(w, "user not found in context", http.StatusUnauthorized)
		return
	}

	if err = h.userService.Deactivate(userID, changedBy); err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, service.ErrUserAlreadyInactive) {
			http.Error(w, "user already inactive", http.StatusBadRequest)
			return
		}
		http.Error(w, "failed to deactivate user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
