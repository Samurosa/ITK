package app

import (
	"ITK_Code/m/v2/app/grps"
	"ITK_Code/m/v2/config"
	"ITK_Code/m/v2/internal/adapters/jwt"
	"ITK_Code/m/v2/internal/adapters/storage/inmemory"
	"ITK_Code/m/v2/internal/application"
	"ITK_Code/m/v2/internal/core/auth"

	"go.uber.org/zap"
)

type App struct {
	GrpcApp *grpsApp.App
}

func New(
	log *zap.Logger,
	port int,
	tokenTTL config.TokensTTL,
	secret string,
) *App {

	userStorage := inmemory.NewUserStorage()
	sessionStorage := inmemory.NewSessionStorage()
	walletStorage := inmemory.NewBalanceStorage()

	secretAuthorization := jwt.NewJWT(auth.JWTConfig{
		Secret:          secret,
		AccessTokenTTL:  tokenTTL.AccessTokenTTL,
		RefreshTokenTTL: tokenTTL.RefreshTokenTTL,
	})

	appService := application.New(log, secretAuthorization, sessionStorage, userStorage, userStorage, walletStorage)

	app := grpsApp.New(log, appService, appService, appService, port)
	return &App{
		GrpcApp: app,
	}
}
