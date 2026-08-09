package app

import (
	"ITK_Code/m/v2/internal/adapters/outbound/repository/postgres"
	"ITK_Code/m/v2/internal/adapters/outbound/repository/redis"
	"ITK_Code/m/v2/internal/config"
	"ITK_Code/m/v2/internal/infrastructure"
	"context"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

type App struct {
	logger *zap.Logger

	ctx      context.Context
	cancel   context.CancelFunc
	postgres *postgres.Storage
	redis    *redis.Storage
	grpcApp  *infrastructure.GRPCApp
}

func New(
	cfg *config.Config,
	secret string,
) (*App, error) {

	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	log := logger.Named("app start")

	ctx, cancel := context.WithCancel(context.Background())

	postgresStorage, err := postgres.NewStorage(ctx,
		log,
		cfg.Postgres.Link,
		cfg.Postgres.MaxRetries,
	)
	if err != nil {
		log.Error("Error starting postgres storage", zap.Error(err))
		cancel()
		return nil, err
	}
	log.Debug("Successfully started postgres storage")

	log.Named("Redis infrastructure")

	redisClient, err := redis.NewRedisClient(ctx, log, cfg.Redis)
	if err != nil {
		log.Error("Error creating redis client", zap.Error(err))
		cancel()
		return nil, err
	}

	redisStorage := redis.NewStorage(redisClient)
	log.Debug("infrastructure initialized")

	services := NewServices(logger, cfg, postgresStorage, redisStorage, secret)

	grpcApp := infrastructure.NewGRPC(logger,
		services,
		cfg.GRPC.Port,
	)

	return &App{
		logger: logger,

		ctx:    ctx,
		cancel: cancel,

		postgres: postgresStorage,
		redis:    redisStorage,

		grpcApp: grpcApp,
	}, nil
}

func (app *App) Start() {
	log := app.logger.Named("starting grpc goroutine")
	go func() {
		err := app.grpcApp.Run()
		if err != nil {
			log.Error("grpc goroutine failed", zap.Error(err))
		}
	}()
}

func (app *App) WaitSignal() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
}

func (app *App) Stop() {
	app.logger.Debug("application stop")

	app.grpcApp.Stop()
	app.cancel()
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
