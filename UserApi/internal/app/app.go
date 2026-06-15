package app

import (
	grpsApp "ITK_Code/m/v2/internal/app/grps"
	"ITK_Code/m/v2/internal/secret/authorization"
	"ITK_Code/m/v2/internal/services"
	"ITK_Code/m/v2/internal/storage/inmemory"
	"time"

	"go.uber.org/zap"
)

type App struct {
	GrpcApp *grpsApp.App
}

func New(
	log *zap.Logger,
	port int,
	tokenTTL time.Duration,
	secret string,
) *App {

	userStorage := inmemory.NewUserStorage()

	secretAuthorization := authorization.NewSecret(secret)

	userService := services.New(log, userStorage, userStorage, secretAuthorization, tokenTTL)

	app := grpsApp.New(log, userService, port)
	return &App{
		GrpcApp: app,
	}
}
