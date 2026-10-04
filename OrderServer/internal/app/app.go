package app

import (
	"ITK_Code/m/v2/internal/adapters/outbound/postgres"
	"ITK_Code/m/v2/internal/adapters/outbound/sessionValidator"
	"ITK_Code/m/v2/internal/adapters/outbound/spot"
	"ITK_Code/m/v2/internal/application"
	"ITK_Code/m/v2/internal/config"
	"ITK_Code/m/v2/internal/infrastructure"
	"context"
	"errors"
	"fmt"
	"os"
	"syscall"

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

func New(cfg *config.Config, secret string) (_ *App, initErr error) {
	log, err := zap.NewProduction()
	if err != nil {
		return nil, fmt.Errorf("initialize logger: %w", err)
	}
	log = log.Named("order_server")
	defer func() {
		if initErr != nil {
			log.Error("order server initialization failed", zap.Error(initErr))
			_ = log.Sync()
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())

	storagePostgres, err := postgres.NewStorage(ctx, log, cfg.Postgres)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("initialize postgres: %w", err)
	}
	log.Info("postgres connected")
	redisClient, err := sessionValidator.NewRedisClient(ctx, cfg.Redis)
	if err != nil {
		storagePostgres.ClosePool()
		cancel()
		return nil, fmt.Errorf("initialize redis: %w", err)
	}
	log.Info("redis connected")

	storageSessions := sessionValidator.NewStorage(redisClient)

	conn, err := grpc.NewClient(
		cfg.Spot.GRPCAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		storagePostgres.ClosePool()
		if closeErr := redisClient.Close(); closeErr != nil {
			log.Warn("redis client cleanup failed", zap.Error(closeErr))
		}
		cancel()
		return nil, fmt.Errorf("initialize spot grpc client: %w", err)
	}

	conn.Connect()
	log.Debug("spot grpc connection requested", zap.String("target", cfg.Spot.GRPCAddr))

	spotClient := spot.NewClient(conn)

	orderService := application.NewOrderService(storagePostgres, spotClient)

	parser, err := jwt.NewParser(secret)
	if err != nil {
		if connErr := conn.Close(); connErr != nil {
			log.Warn("spot grpc client cleanup failed", zap.Error(connErr))
		}
		storagePostgres.ClosePool()
		if closeErr := redisClient.Close(); closeErr != nil {
			log.Warn("redis client cleanup failed", zap.Error(closeErr))
		}
		cancel()
		return nil, fmt.Errorf("initialize jwt parser: %w", err)
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
	log := app.logger.Named("grpc")
	go func() {
		err := app.grpcApp.Run()
		if err != nil {
			log.Error("grpc server failed", zap.Error(err))
		}
	}()
}

func (app *App) Stop() {
	app.logger.Info("order server stopping")

	app.grpcApp.Stop()
	app.cancel()
	if err := app.sessions.Close(); err != nil {
		app.logger.Warn("redis session client close failed", zap.Error(err))
	}
	app.postgres.ClosePool()
	if err := app.spotConn.Close(); err != nil {
		app.logger.Warn("spot grpc client close failed", zap.Error(err))
	}
	app.logger.Info("order server stopped")

	err := app.logger.Sync()
	if err != nil && !errors.Is(err, syscall.EINVAL) && !errors.Is(err, syscall.ENOTTY) {
		// A failed sink cannot reliably report its own flush failure.
		_, _ = fmt.Fprintf(os.Stderr, "flush order server logger: %v\n", err)
	}
}
