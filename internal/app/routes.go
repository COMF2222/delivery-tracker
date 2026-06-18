package app

import (
	"delivery-tracker/internal/domain"
	"delivery-tracker/internal/handler"
	"net/http"
)

func RegisterRoutes(deps *Dependencies) {
	protect := func(handler http.HandlerFunc, roles ...domain.Role) http.HandlerFunc {
		return deps.AuthMiddleware.RequireAuth(
			deps.AuthMiddleware.RequireRole(roles...)(handler),
		)
	}
	http.HandleFunc("/health", handler.Health)

	http.HandleFunc("/api/v1/parcels",
		protect(deps.ParcelHandler.CreateParcel, domain.RoleAdmin, domain.RoleManager))
	http.HandleFunc("/api/v1/parcels/status",
		protect(deps.ParcelHandler.UpdateStatus, domain.RoleAdmin, domain.RoleManager))
	http.HandleFunc("/api/v1/parcels/photos",
		protect(deps.ParcelHandler.AddPhoto, domain.RoleAdmin, domain.RoleManager))
	http.HandleFunc("/api/v1/parcels/archive",
		protect(deps.ParcelHandler.Archive, domain.RoleAdmin, domain.RoleManager))
	http.HandleFunc("/api/v1/parcels/track", deps.ParcelHandler.GetByTrackNumber)

	http.HandleFunc("/api/v1/users", protect(deps.UserHandler.Create, domain.RoleAdmin))

	http.HandleFunc("/api/v1/auth/login", deps.AuthHandler.Login)
}
