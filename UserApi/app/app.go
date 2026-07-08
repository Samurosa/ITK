package app

import (
	"ITK_Code/m/v2/app/grps"
	"ITK_Code/m/v2/app/workers"
	"ITK_Code/m/v2/config"
	"ITK_Code/m/v2/internal/adapters/jwt"
	"ITK_Code/m/v2/internal/adapters/storage/postgres"
	"ITK_Code/m/v2/internal/application"
	"ITK_Code/m/v2/internal/core/auth"
	"context"

	"go.uber.org/zap"
)

type App struct {
	GrpcApp *grpsApp.App
	Workers *workers.App
}

func New(
	log *zap.Logger,
	storagePath string,
	port int,
	tokenTTL config.TokensTTL,
	secret string,
) *App {
	//TODO: переделать сделать cleaner для контекста и соединения с бд
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	postgresStorage, err := postgres.NewPool(ctx, log, storagePath, 10)
	if err != nil {
		panic(err)
	}
	pool := postgresStorage.GetPool()
	defer pool.Close()

	userStorage := postgres.NewUserStorage(pool)

	sessionStorage := postgres.NewSessionStorage(pool)

	walletStorage := postgres.NewBalanceStorage(pool)

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

	workersApp := workers.NewWorker(log, ctx, cancel, sessionStorage)

	return &App{
		GrpcApp: app,
		Workers: workersApp,
	}
}
