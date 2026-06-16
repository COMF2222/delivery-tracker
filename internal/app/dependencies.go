package app

import (
	"delivery-tracker/internal/handler"
	"delivery-tracker/internal/repository"
	"delivery-tracker/internal/service"

	"github.com/jmoiron/sqlx"
)

type Dependencies struct {
	ParcelHandler *handler.ParcelHandler
	UserHandler   *handler.UserHandler
	AuthHandler   *handler.AuthHandler
}

func NewDependencies(db *sqlx.DB) *Dependencies {
	parcelRepo := repository.NewParcelRepository(db)
	statusRepo := repository.NewStatusRepository(db)
	photoRepo := repository.NewParcelPhotoRepository(db)
	historyRepo := repository.NewParcelStatusHistoryRepository(db)
	auditRepo := repository.NewAuditRepository(db)
	userRepo := repository.NewUserRepository(db)

	txManager := repository.NewTransactionManager(db)

	parcelService := service.NewParcelService(
		parcelRepo,
		statusRepo,
		photoRepo,
		historyRepo,
		auditRepo,
		txManager,
	)

	userService := service.NewUserService(userRepo)

	authService := service.NewAuthService(userRepo)

	return &Dependencies{
		ParcelHandler: handler.NewParcelHandler(parcelService),
		UserHandler:   handler.NewUserHandler(userService),
		AuthHandler:   handler.NewAuthHandler(authService),
	}
}
