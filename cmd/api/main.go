package main

import (
	"delivery-tracker/internal/config"
	"delivery-tracker/internal/database"
	"delivery-tracker/internal/handler"
	"delivery-tracker/internal/repository"
	"delivery-tracker/internal/service"
	"fmt"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("config loaded")
	db, err := database.New(cfg.Database)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("failed to close database: %v", err)
		}
	}()

	fmt.Println("database connected")
	parcelRepo := repository.NewParcelRepository(db)
	statusRepo := repository.NewStatusRepository(db)
	photoRepo := repository.NewParcelPhotoRepository(db)
	historyRepo := repository.NewParcelStatusHistoryRepository(db)
	auditRepo := repository.NewAuditRepository(db)
	txManger := repository.NewTransactionManager(db)

	parcelService := service.NewParcelService(parcelRepo, statusRepo, photoRepo, historyRepo, auditRepo, txManger)

	health := handler.Health
	parcelHandler := handler.NewParcelHandler(parcelService)

	http.HandleFunc("/health", health)
	http.HandleFunc("/api/v1/parcels", parcelHandler.CreateParcel)
	http.HandleFunc("/api/v1/parcels/track", parcelHandler.GetByTrackNumber)
	http.HandleFunc("/api/v1/parcels/status", parcelHandler.UpdateStatus)
	http.HandleFunc("/api/v1/parcels/photos", parcelHandler.AddPhoto)

	_ = http.ListenAndServe(":"+cfg.Server.Port, nil)
}
