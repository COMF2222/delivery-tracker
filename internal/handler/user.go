package handler

import (
	"delivery-tracker/internal/contextkeys"
	"delivery-tracker/internal/dto"
	"delivery-tracker/internal/repository"
	"delivery-tracker/internal/response"
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

// Create создает пользователя.
//
//	@Summary		Создание пользователя
//	@Description	Создаёт пользователя
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.CreateUserRequest	true	"Create user request"
//	@Success		201		{object}	dto.CreateUserResponse
//	@Failure		400		{object}	response.ErrorResponse	"Bad request"
//	@Failure		401		{object}	response.ErrorResponse	"Unauthorized"
//	@Failure		403		{object}	response.ErrorResponse	"Forbidden"
//	@Failure		405		{object}	response.ErrorResponse	"Method not allowed"
//	@Failure		500		{object}	response.ErrorResponse	"Internal server error"
//	@Security		BearerAuth
//	@Router			/users [post]
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		response.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req dto.CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "failed to decode request", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.userService.CreateUser(req.Login, req.Password, req.Role)
	if err != nil {
		response.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	resp := dto.CreateUserResponse{ID: user.ID}

	response.JSON(w, http.StatusCreated, resp)
}

// Deactivate деактивирует пользователя по ID.
//
//	@Summary		Деактивация пользователя
//	@Description	Деактивирует пользователя по ID
//	@Tags			Users
//	@Produce		json
//	@Param			id	query	int	true	"User ID"
//	@Success		204	"No Content"
//	@Failure		400	{object}	response.ErrorResponse	"Bad request"
//	@Failure		401	{object}	response.ErrorResponse	"Unauthorized"
//	@Failure		403	{object}	response.ErrorResponse	"Forbidden"
//	@Failure		404	{object}	response.ErrorResponse	"Not found"
//	@Failure		405	{object}	response.ErrorResponse	"Method not allowed"
//	@Failure		500	{object}	response.ErrorResponse	"Internal server error"
//	@Security		BearerAuth
//	@Router			/users/deactivate [patch]
func (h *UserHandler) Deactivate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PATCH" {
		response.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		response.Error(w, "id cannot be empty", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(id)
	if err != nil {
		response.Error(w, "failed convert id to int", http.StatusBadRequest)
		return
	}
	if userID <= 0 {
		response.Error(w, "id must be positive", http.StatusBadRequest)
		return
	}

	changedBy, ok := r.Context().Value(contextkeys.UserID).(int)
	if !ok {
		response.Error(w, "user not found in context", http.StatusUnauthorized)
		return
	}

	if err = h.userService.Deactivate(userID, changedBy); err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			response.Error(w, "user not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, service.ErrUserAlreadyInactive) {
			response.Error(w, "user already inactive", http.StatusBadRequest)
			return
		}
		response.Error(w, "failed to deactivate user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
