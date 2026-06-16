package app

import (
	"delivery-tracker/internal/config"
	"delivery-tracker/internal/database"
	"log"
	"net/http"
)

func Run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	db, err := database.New(cfg.Database)
	if err != nil {
		return err
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("failed to close database: %v", err)
		}
	}()

	deps := NewDependencies(db, &cfg)

	RegisterRoutes(deps)

	return http.ListenAndServe(":"+cfg.Server.Port, nil)
}
