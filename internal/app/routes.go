package app

import (
	"delivery-tracker/internal/handler"
	"net/http"
)

func RegisterRoutes(deps *Dependencies) {
	http.HandleFunc("/health", handler.Health)

	http.HandleFunc("/api/v1/parcels", deps.AuthMiddleware.RequireAuth(deps.ParcelHandler.CreateParcel))
	http.HandleFunc("/api/v1/parcels/track", deps.ParcelHandler.GetByTrackNumber)
	http.HandleFunc("/api/v1/parcels/status", deps.AuthMiddleware.RequireAuth(deps.ParcelHandler.UpdateStatus))
	http.HandleFunc("/api/v1/parcels/photos", deps.AuthMiddleware.RequireAuth(deps.ParcelHandler.AddPhoto))

	http.HandleFunc("/api/v1/users", deps.UserHandler.Create)

	http.HandleFunc("/api/v1/auth/login", deps.AuthHandler.Login)
}
