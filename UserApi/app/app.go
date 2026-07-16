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

	postgres *postgres.Storage

	redis *redis.Storage
}

func New(
	cfg *config.Config,
	secret string,
) (*App, error) {

	log, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	ctxApp := contextStatus.New()
	ctx := ctxApp.GetContext()

	postgresStorage, err := postgres.NewPool(ctx, log, cfg.Storage.Link, 10)
	if err != nil {
		return nil, err
	}

	pool := postgresStorage.GetPool()

	userStorage := postgres.NewUserStorage(pool)

	redisClient, err := redis.NewRedisClient(ctx, cfg.Redis)
	if err != nil {
		return nil, err
	}

	sessionStorage := redis.NewStorage(redisClient)

	walletStorage := postgres.NewBalanceStorage(pool)

	secretAuthorization := jwt.NewJWT(log,
		auth.JWTConfig{
			Secret:          secret,
			AccessTokenTTL:  cfg.TokenTTl.AccessTokenTTL,
			RefreshTokenTTL: cfg.TokenTTl.RefreshTokenTTL,
		})

	userService := application.NewUserService(log,
		userStorage,
		userStorage,
	)

	authService := application.NewAuthService(log,
		secretAuthorization,
		sessionStorage,
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
		cfg.GRPC.Port,
	)

	workersApp := workers.NewWorker(log, ctx, sessionStorage)

	return &App{
		logger: log,

		grpcApp: app,

		workers: workersApp,

		context: ctxApp,

		postgres: postgresStorage,

		redis: sessionStorage,
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
	app.postgres.ClosePool()
	err := app.redis.Stop()
	if err != nil {
		app.logger.Error("redis stop", zap.Error(err))
	}
	err = app.logger.Sync()
	if err != nil {
		app.logger.Error("error sync logger: ", zap.Error(err))
	}
}
