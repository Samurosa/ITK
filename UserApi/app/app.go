package app

import (
	"ITK_Code/m/v2/app/grps"
	"ITK_Code/m/v2/internal/adapters/storage/inmemory"
	"ITK_Code/m/v2/internal/application/secret/authorization"
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

	userService := user.New(log, userStorage, userStorage, secretAuthorization, tokenTTL)

	app := grpsApp.New(log, userService, port)
	return &App{
		GrpcApp: app,
	}
}
