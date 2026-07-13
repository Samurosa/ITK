package app

import (
	"ITK_Code/m/v2/app/contextStatus"
	"ITK_Code/m/v2/app/grps"
	"ITK_Code/m/v2/app/workers"
	"ITK_Code/m/v2/config"
	"ITK_Code/m/v2/internal/adapters/jwt"
	"ITK_Code/m/v2/internal/adapters/storage/postgres"
	"ITK_Code/m/v2/internal/adapters/storage/redis"
	"ITK_Code/m/v2/internal/application"
	"ITK_Code/m/v2/internal/core/auth"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

type App struct {
	logger *zap.Logger

	grpcApp *grpsApp.App

	workers *workers.App

	context *contextStatus.App

	storage *postgres.Storage
}

func New(
	storagePath string,
	port int,
	tokenTTL config.TokensTTL,
	secret string,
) (*App, error) {

	log, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	ctxApp := contextStatus.New()
	ctx := ctxApp.GetContext()

	postgresStorage, err := postgres.NewPool(ctx, log, storagePath, 10)
	if err != nil {
		return nil, err
	}

	pool := postgresStorage.GetPool()

	userStorage := postgres.NewUserStorage(pool)

	sessionStorage := redis.NewSessionStorage(pool)

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
		sessionStorage,
		port,
	)

	workersApp := workers.NewWorker(log, ctx, sessionStorage)

	return &App{
		logger: log,

		grpcApp: app,

		workers: workersApp,

		context: ctxApp,

		storage: postgresStorage,
	}, nil
}

func (app *App) Start() {
	go app.grpcApp.Run()
	go app.workers.Run()
}

func (app *App) WaitSignal() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
}

func (app *App) Stop() {
	app.logger.Debug("application stop")

	app.grpcApp.Stop()
	app.context.Stop()
	app.storage.ClosePool()
	err := app.logger.Sync()
	if err != nil {
		app.logger.Error("error sync logger: ", zap.Error(err))
	}
}
