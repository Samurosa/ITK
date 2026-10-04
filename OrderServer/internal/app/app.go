package app

import (
	"ITK_Code/m/v2/internal/adapters/outbound/postgres"
	"ITK_Code/m/v2/internal/adapters/outbound/sessionValidator"
	"ITK_Code/m/v2/internal/adapters/outbound/spot"
	"ITK_Code/m/v2/internal/application"
	"ITK_Code/m/v2/internal/config"
	"ITK_Code/m/v2/internal/infrastructure"
	"context"
	"fmt"

	"github.com/Samurosa/exchange-common/shared/auth/jwt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	logger *zap.Logger

	ctx    context.Context
	cancel context.CancelFunc

	postgres *postgres.Storage
	sessions *sessionValidator.Storage
	spotConn *grpc.ClientConn

	grpcApp *infrastructure.GRPCApp
}

func New(cfg *config.Config, secret string) (*App, error) {
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
	redisClient, err := sessionValidator.NewRedisClient(ctx, log, cfg.Redis)
	if err != nil {
		storagePostgres.ClosePool()
		cancel()
		log.Error("Failed to connect to redis", zap.Error(err))
		return nil, err
	}

	storageSessions := sessionValidator.NewStorage(redisClient)

	conn, err := grpc.NewClient(
		cfg.Spot.GRPCAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		storagePostgres.ClosePool()
		_ = redisClient.Close()
		cancel()
		return nil, err
	}

	conn.Connect()

	spotClient := spot.NewClient(conn)

	orderService := application.NewOrderService(log, storagePostgres, spotClient)

	parser, err := jwt.NewParser(secret)
	if err != nil {
		if connErr := conn.Close(); connErr != nil {
			log.Error("Failed to close connection", zap.Error(connErr))
		}
		storagePostgres.ClosePool()
		_ = redisClient.Close()
		cancel()
		return nil, err
	}

	grpcServer := infrastructure.NewGRPC(log, orderService, parser, storageSessions, cfg.GRPC.Port)

	return &App{
		logger: log,
		ctx:    ctx,
		cancel: cancel,

		postgres: storagePostgres,
		sessions: storageSessions,
		spotConn: conn,

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
	if err := app.sessions.Close(); err != nil {
		app.logger.Error("close redis session client", zap.Error(err))
	}
	app.postgres.ClosePool()
	if err := app.spotConn.Close(); err != nil {
		app.logger.Error("close spot client connection", zap.Error(err))
	}

	err := app.logger.Sync()
	if err != nil {
		app.logger.Error("error sync logger: ", zap.Error(err))
	}
}
