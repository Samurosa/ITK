package app

import (
	"ITK_Code/m/v2/internal/adapters/outbound/postgres"
	"ITK_Code/m/v2/internal/application"
	"ITK_Code/m/v2/internal/config"
	"ITK_Code/m/v2/internal/infrastructure"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

type App struct {
	logger *zap.Logger

	ctx    context.Context
	cancel context.CancelFunc

	postgres *postgres.Storage

	grpcApp *infrastructure.GRPCApp
}

func New(cfg *config.Config) (*App, error) {
	log, err := zap.NewProduction()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	storagePostgres, err := postgres.NewStorage(ctx, log, cfg.Postgres)
	if err != nil {
		cancel()
		log.Error("Failed to connect to postgres", zap.Error(err))
		return nil, err
	}

	spotService := application.NewSpot(log, storagePostgres)

	grpcServer := infrastructure.NewGRPC(log, spotService, cfg.GRPC.Port)

	return &App{
		logger: log,
		ctx:    ctx,
		cancel: cancel,

		postgres: storagePostgres,

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

	err := app.logger.Sync()
	if err != nil {
		app.logger.Error("error sync logger: ", zap.Error(err))
	}
}
