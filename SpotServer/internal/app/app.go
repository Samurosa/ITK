package app

import (
	"ITK_Code/m/v2/internal/adapters/outbound/postgres"
	"ITK_Code/m/v2/internal/adapters/outbound/sessionValidator"
	"ITK_Code/m/v2/internal/application"
	"ITK_Code/m/v2/internal/config"
	"ITK_Code/m/v2/internal/infrastructure"
	"context"
	"fmt"

	"github.com/Samurosa/exchange-common/shared/auth/jwt"

	"go.uber.org/zap"
)

type App struct {
	logger *zap.Logger

	ctx    context.Context
	cancel context.CancelFunc

	postgres *postgres.Storage
	redis    *sessionValidator.Storage

	grpcApp *infrastructure.GRPCApp
}

func New(cfg *config.Config, secret string) (*App, error) {
	log, err := zap.NewProduction()
	if err != nil {
		fmt.Println("failed to initialize logger")
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	parser, err := jwt.NewParser(secret)
	if err != nil {
		cancel()
		return nil, err
	}

	storagePostgres, err := postgres.NewStorage(ctx, log, cfg.Postgres)
	if err != nil {
		cancel()
		log.Error("Failed to connect to postgres", zap.Error(err))
		return nil, err
	}

	spotService := application.NewSpot(log, storagePostgres)
	redisStorage, err := sessionValidator.NewStorage(ctx, cfg.Redis)
	if err != nil {
		storagePostgres.ClosePool()
		cancel()
		return nil, err
	}

	grpcServer := infrastructure.NewGRPC(log, spotService, cfg.GRPC.Port, parser, redisStorage)

	return &App{
		logger: log,
		ctx:    ctx,
		cancel: cancel,

		postgres: storagePostgres,
		redis:    redisStorage,

		grpcApp: grpcServer,
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

func (app *App) Stop() {
	app.logger.Debug("application stop")

	app.grpcApp.Stop()
	app.cancel()
	app.postgres.ClosePool()
	if err := app.redis.Close(); err != nil {
		app.logger.Error("redis close", zap.Error(err))
	}

	err := app.logger.Sync()
	if err != nil {
		app.logger.Error("error sync logger: ", zap.Error(err))
	}
}
