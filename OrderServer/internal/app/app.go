package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"ITK_Code/m/v2/internal/adapters/outbound/postgres"
	"ITK_Code/m/v2/internal/adapters/outbound/sessionValidator"
	"ITK_Code/m/v2/internal/adapters/outbound/spot"
	"ITK_Code/m/v2/internal/application"
	"ITK_Code/m/v2/internal/config"
	"ITK_Code/m/v2/internal/infrastructure"

	"github.com/Samurosa/exchange-common/shared/auth/jwt"
	sharedredis "github.com/Samurosa/exchange-common/shared/redis"
	"github.com/Samurosa/exchange-common/shared/telemetry"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	logger     *zap.Logger
	postgres   *postgres.Storage
	sessions   *sessionValidator.Storage
	spotConn   *grpc.ClientConn
	grpcApp    *infrastructure.GRPCApp
	shutDowner fx.Shutdowner
}

func Run(configPath string, logger *zap.Logger) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf(
			"load config from %s: %w",
			configPath,
			err,
		)
	}

	secret, err := os.ReadFile(cfg.JWTSecretPath)
	if err != nil {
		return fmt.Errorf(
			"read JWT secret from %s: %w",
			cfg.JWTSecretPath,
			err,
		)
	}

	shutdownTracing, err := telemetry.InitTracing(context.Background(), "order-service")
	if err != nil {
		return fmt.Errorf(
			"initialize tracing: %w",
			err,
		)
	}

	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := shutdownTracing(shutdownCtx); err != nil {
			logger.Error(
				"failed to shutdown tracing",
				zap.Error(err),
			)
		}
	}()

	fxApp := fx.New(
		fx.Supply(
			cfg,
			string(secret),
			logger,
		),
		fx.Provide(New),
		fx.Invoke(registerLifecycle),
		fx.StartTimeout(30*time.Second),
		fx.StopTimeout(15*time.Second),
	)

	if err := fxApp.Err(); err != nil {
		return fmt.Errorf(
			"initialize application: %w",
			err,
		)
	}

	fxApp.Run()

	return nil
}

func registerLifecycle(
	lifecycle fx.Lifecycle,
	application *App,
) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err := application.Start(ctx); err != nil {
				cleanupErr := application.closeResources()

				return errors.Join(err, cleanupErr)
			}

			return nil
		},

		OnStop: func(ctx context.Context) error {
			return application.Stop(ctx)
		},
	})
}

func New(
	cfg *config.Config,
	secret string,
	logger *zap.Logger,
	shutDowner fx.Shutdowner,
) (*App, error) {
	orderApp := &App{
		logger:     logger.Named("order_server"),
		shutDowner: shutDowner,
	}

	initialized := false

	defer func() {
		if !initialized {
			if err := orderApp.closeResources(); err != nil {
				orderApp.logger.Error(
					"failed to clean up resources",
					zap.Error(err),
				)
			}
		}
	}()

	log := orderApp.logger

	initCtx, cancel := context.WithTimeout(
		context.Background(),
		30*time.Second,
	)
	defer cancel()

	parser, err := jwt.NewParser(secret)
	if err != nil {
		return nil, fmt.Errorf(
			"initialize jwt parser: %w",
			err,
		)
	}

	storagePostgres, err := postgres.NewStorage(
		initCtx,
		log,
		cfg.Postgres,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"initialize postgres: %w",
			err,
		)
	}

	orderApp.postgres = storagePostgres

	log.Info("postgres connected")

	redisClient, err := sharedredis.NewClient(
		initCtx,
		log,
		cfg.Redis,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"initialize redis: %w",
			err,
		)
	}

	orderApp.sessions = sessionValidator.NewStorage(
		redisClient,
	)

	log.Info("redis connected")

	conn, err := grpc.NewClient(
		cfg.Spot.GRPCAddr,
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
		grpc.WithStatsHandler(
			otelgrpc.NewClientHandler(),
		),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"initialize spot grpc client: %w",
			err,
		)
	}

	orderApp.spotConn = conn

	conn.Connect()

	log.Debug(
		"spot grpc connection requested",
		zap.String("target", cfg.Spot.GRPCAddr),
	)

	spotClient := spot.NewClient(conn)

	orderService := application.NewOrderService(
		storagePostgres,
		spotClient,
	)

	orderApp.grpcApp = infrastructure.NewGRPC(
		logger.Named("grpc"),
		orderService,
		parser,
		orderApp.sessions,
		cfg.GRPC.Port,
	)

	initialized = true

	return orderApp, nil
}

func (a *App) Start(ctx context.Context) error {
	err := a.grpcApp.Start(
		ctx,
		func(serveErr error) {
			a.logger.Error(
				"grpc server terminated unexpectedly",
				zap.Error(serveErr),
			)

			if shutdownErr := a.shutDowner.Shutdown(
				fx.ExitCode(1),
			); shutdownErr != nil {
				a.logger.Error(
					"failed to request application shutdown",
					zap.Error(shutdownErr),
				)
			}
		},
	)

	if err != nil {
		return fmt.Errorf(
			"start application: %w",
			err,
		)
	}

	return nil
}

func (a *App) Stop(ctx context.Context) error {
	a.logger.Info("order server stopping")

	var grpcErr error

	if a.grpcApp != nil {
		grpcErr = a.grpcApp.Stop(ctx)
	}

	cleanupErr := a.closeResources()

	a.logger.Info("order server stopped")

	return errors.Join(grpcErr, cleanupErr)
}

func (a *App) closeResources() error {
	var errs []error

	if a.spotConn != nil {
		if err := a.spotConn.Close(); err != nil {
			a.logger.Warn(
				"spot grpc client close failed",
				zap.Error(err),
			)

			errs = append(errs, err)
		} else {
			a.logger.Info("spot grpc client closed")
		}
	}

	if a.sessions != nil {
		if err := a.sessions.Close(); err != nil {
			a.logger.Warn(
				"redis session client close failed",
				zap.Error(err),
			)

			errs = append(errs, err)
		} else {
			a.logger.Info("redis connection closed")
		}
	}

	if a.postgres != nil {
		a.postgres.ClosePool()
		a.logger.Info("postgres pool closed")
	}

	return errors.Join(errs...)
}
