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

	secretAuthorization := jwt.NewJWT(log,
		auth.JWTConfig{
			Secret:          secret,
			AccessTokenTTL:  tokenTTL.AccessTokenTTL,
			RefreshTokenTTL: tokenTTL.RefreshTokenTTL,
		})

	userService := application.NewUserService(log,
		userStorage,
		userStorage,
	)

	authService := application.NewAuthService(log,
		secretAuthorization,
		sessionStorage,
		userStorage,
		userStorage,
	)

	walletService := application.NewWalletService(log,
		walletStorage,
	)

	app := grpsApp.New(log,
		secretAuthorization,
		userService,
		authService,
		walletService,
		port,
	)
	return &App{
		GrpcApp: app,
	}
}
