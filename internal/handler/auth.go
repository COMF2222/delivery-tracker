package handler

import (
	"delivery-tracker/internal/dto"
	"delivery-tracker/internal/repository"
	"delivery-tracker/internal/response"
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

// Login
//
//	@Summary		Авторизация пользователя
//	@Description	Авторизация пользователей
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.LoginUserRequest	true	"Login request"
//	@Success		200		{object}	dto.LoginUserResponse
//	@Failure		400		{object}	response.ErrorResponse	"Bad request"
//	@Failure		401		{object}	response.ErrorResponse	"Invalid login or password"
//	@Failure		403		{object}	response.ErrorResponse	"User inactive"
//	@Failure		405		{object}	response.ErrorResponse	"Method not allowed"
//	@Failure		500		{object}	response.ErrorResponse	"Internal server error"
//	@Router			/auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req dto.LoginUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "failed to decode request", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := h.authService.Login(req.Login, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidPassword) || errors.Is(err, repository.ErrUserNotFound) {
			response.Error(w, "invalid login or password", http.StatusUnauthorized)
			return
		}
		if errors.Is(err, service.ErrUserInactive) {
			response.Error(w, "user is inactive", http.StatusForbidden)
			return
		}
		response.Error(w, "failed to login user", http.StatusInternalServerError)
		return
	}

	resp := dto.LoginUserResponse{Token: token}
	response.JSON(w, http.StatusOK, resp)
}
