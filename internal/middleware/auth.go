package middleware

import (
	"context"
	"delivery-tracker/internal/contextkeys"
	"delivery-tracker/internal/domain"
	"delivery-tracker/internal/service"
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	authService *service.AuthService
}

func NewAuthMiddleware(authService *service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{authService: authService}
}

func (m *AuthMiddleware) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			http.Error(w, "invalid token format", http.StatusUnauthorized)
			return
		}

		userID, login, role, err := m.authService.ValidateToken(token)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, contextkeys.UserID, userID)
		ctx = context.WithValue(ctx, contextkeys.Login, login)
		ctx = context.WithValue(ctx, contextkeys.Role, role)

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	}
}

func (m *AuthMiddleware) RequireRole(allowedRoles ...domain.Role) func(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			userRole, ok := r.Context().Value(contextkeys.Role).(domain.Role)
			if !ok {
				http.Error(w, "the user is not authenticated", http.StatusUnauthorized)
				return
			}

			for _, allowedRole := range allowedRoles {
				if userRole == allowedRole {
					next.ServeHTTP(w, r)
					return
				}
			}

			http.Error(w, "insufficient access rights", http.StatusForbidden)
			return
		}
	}
}
