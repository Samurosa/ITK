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

	log, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	log.Named("app start")

	contextModel := infrastructure.NewContext()
	ctx := contextModel.GetContext()

	storages, err := infrastructure.NewStorages(ctx, log, cfg)
	if err != nil {
		log.Error("failed to initialize infrastructure", zap.Error(err))
		return nil, err
	}
	log.Info("infrastructure initialized")

	services := NewServices(log, cfg, storages, secret)

	grpcApp := infrastructure.NewGRPC(log,
		services,
		cfg.GRPC.Port,
	)

	return &App{
		logger: log,

		context: contextModel,

		storages: storages,

		grpcApp: grpcApp,
	}, nil
}

func (app *App) Start() {
	go app.grpcApp.Run()
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
