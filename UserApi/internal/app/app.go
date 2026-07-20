package app

import (
	"ITK_Code/m/v2/config"
	"ITK_Code/m/v2/internal/infrastructure"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

type App struct {
	logger *zap.Logger

	context  *infrastructure.Context
	storages *infrastructure.Storages
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

	contextModel := infrastructure.NewContext()
	ctx := contextModel.GetContext()

	storages, err := infrastructure.NewStorages(ctx, logger, cfg)
	if err != nil {
		log.Error("failed to initialize infrastructure", zap.Error(err))
		return nil, err
	}
	log.Info("infrastructure initialized")

	services := NewServices(logger, cfg, storages, secret)

	grpcApp := infrastructure.NewGRPC(logger,
		services,
		cfg.GRPC.Port,
	)

	return &App{
		logger: logger,

		context: contextModel,

		storages: storages,

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
	app.context.Stop()
	app.storages.Postgres.ClosePool()
	err := app.storages.Redis.Stop()
	if err != nil {
		app.logger.Error("redis stop", zap.Error(err))
	}
	err = app.logger.Sync()
	if err != nil {
		app.logger.Error("error sync logger: ", zap.Error(err))
	}
}
