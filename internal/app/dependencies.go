package app

import (
	"delivery-tracker/internal/cache"
	"delivery-tracker/internal/config"
	"delivery-tracker/internal/generator"
	"delivery-tracker/internal/handler"
	"delivery-tracker/internal/middleware"
	"delivery-tracker/internal/repository"
	"delivery-tracker/internal/service"
	"github.com/redis/go-redis/v9"

	"github.com/jmoiron/sqlx"
)

type Dependencies struct {
	ParcelHandler  *handler.ParcelHandler
	UserHandler    *handler.UserHandler
	AuthHandler    *handler.AuthHandler
	AuthMiddleware *middleware.AuthMiddleware
}

func NewDependencies(db *sqlx.DB, client *redis.Client, cfg *config.Config) *Dependencies {
	parcelRepo := repository.NewParcelRepository(db)
	statusRepo := repository.NewStatusRepository(db)
	photoRepo := repository.NewParcelPhotoRepository(db)
	historyRepo := repository.NewParcelStatusHistoryRepository(db)
	auditRepo := repository.NewAuditRepository(db)
	userRepo := repository.NewUserRepository(db)
	parcelCache := cache.NewParcelCache(client)

	txManager := repository.NewTransactionManager(db)

	trackGenerator := generator.NewTrackNumberGenerator()

	parcelService := service.NewParcelService(
		parcelRepo,
		statusRepo,
		photoRepo,
		historyRepo,
		auditRepo,
		parcelCache,
		txManager,
		trackGenerator,
	)

	userService := service.NewUserService(userRepo, auditRepo, txManager)

	authService := service.NewAuthService(userRepo, cfg.JWT.Secret, cfg.JWT.TTL)

	authMiddleware := middleware.NewAuthMiddleware(authService)
	return &Dependencies{
		ParcelHandler:  handler.NewParcelHandler(parcelService),
		UserHandler:    handler.NewUserHandler(userService),
		AuthHandler:    handler.NewAuthHandler(authService),
		AuthMiddleware: authMiddleware,
	}
}
