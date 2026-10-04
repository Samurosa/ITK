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
		return nil, fmt.Errorf("initialize logger: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	parser, err := jwt.NewParser(secret)
	if err != nil {
		log.Error("jwt parser initialization failed", zap.Error(err))
		cancel()
		_ = log.Sync()
		return nil, fmt.Errorf("initialize jwt parser: %w", err)
	}

	storagePostgres, err := postgres.NewStorage(ctx, log, cfg.Postgres)
	if err != nil {
		cancel()
		log.Error("postgres initialization failed", zap.Error(err))
		_ = log.Sync()
		return nil, fmt.Errorf("initialize postgres: %w", err)
	}

	spotService := application.NewSpot(storagePostgres)
	redisStorage, err := sessionValidator.NewStorage(ctx, cfg.Redis)
	if err != nil {
		log.Error("redis initialization failed", zap.Error(err))
		storagePostgres.ClosePool()
		cancel()
		_ = log.Sync()
		return nil, fmt.Errorf("initialize redis: %w", err)
	}
	log.Info("redis connected")

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
	log := app.logger.Named("grpc")
	go func() {
		err := app.grpcApp.Run()
		if err != nil {
			log.Error("grpc server failed", zap.Error(err))
		}
	}()
}

func (app *App) Stop() {
	app.logger.Info("application stopping")

	app.grpcApp.Stop()
	app.cancel()
	app.postgres.ClosePool()
	app.logger.Info("postgres pool closed")
	var closeErr error
	if err := app.redis.Close(); err != nil {
		closeErr = err
		app.logger.Error("redis close failed", zap.Error(err))
	} else {
		app.logger.Info("redis connection closed")
	}
	app.logger.Info("application stopped", zap.Bool("cleanup_successful", closeErr == nil))

	err := app.logger.Sync()
	if err != nil {
		app.logger.Warn("logger sync failed", zap.Error(err))
	}
}
