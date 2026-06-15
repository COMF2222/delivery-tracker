package app

import (
	"delivery-tracker/internal/handler"
	"net/http"
)

func RegisterRoutes(deps *Dependencies) {
	http.HandleFunc("/health", handler.Health)

	http.HandleFunc("/api/v1/parcels", deps.ParcelHandler.CreateParcel)
	http.HandleFunc("/api/v1/parcels/track", deps.ParcelHandler.GetByTrackNumber)
	http.HandleFunc("/api/v1/parcels/status", deps.ParcelHandler.UpdateStatus)
	http.HandleFunc("/api/v1/parcels/photos", deps.ParcelHandler.AddPhoto)

	http.HandleFunc("/api/v1/users", deps.UserHandler.Create)
}
