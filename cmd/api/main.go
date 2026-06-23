package main

import (
	_ "delivery-tracker/docs"
	"delivery-tracker/internal/app"
	"log"
)

// @title						Delivery Tracker API
// @version					1.0
// @description				Parcel tracking service
// @host						localhost:8081
// @BasePath					/api/v1
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
