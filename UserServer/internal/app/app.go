package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"ITK_Code/m/v2/internal/adapters/outbound/crypto/hash"
	"ITK_Code/m/v2/internal/adapters/outbound/crypto/jwt"
	"ITK_Code/m/v2/internal/adapters/outbound/repository/postgres"
	"ITK_Code/m/v2/internal/adapters/outbound/repository/redis"
	applicationCore "ITK_Code/m/v2/internal/application"
	"ITK_Code/m/v2/internal/config"
	"ITK_Code/m/v2/internal/infrastructure"

	sharedjwt "github.com/Samurosa/exchange-common/shared/auth/jwt"
	sharedsession "github.com/Samurosa/exchange-common/shared/auth/session"
	sharedredis "github.com/Samurosa/exchange-common/shared/redis"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

type App struct {
	logger     *zap.Logger
	postgres   *postgres.Storage
	redis      *redis.Storage
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
	application := &App{
		logger:     logger.Named("app"),
		shutDowner: shutDowner,
	}

	initialized := false

	defer func() {
		if !initialized {
			if err := application.closeResources(); err != nil {
				application.logger.Error(
					"failed to clean up resources",
					zap.Error(err),
				)
			}
		}
	}()

	log := application.logger

	initCtx, cancel := context.WithTimeout(
		context.Background(),
		30*time.Second,
	)
	defer cancel()

	postgresStorage, err := postgres.NewStorage(
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

	application.postgres = postgresStorage

	log.Info("PostgreSQL connected")

	redisClient, err := sharedredis.NewClient(
		initCtx,
		application.logger,
		cfg.Redis,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"initialize redis: %w",
			err,
		)
	}

	application.redis = redis.NewStorage(redisClient)

	log.Info("Redis connected")

	userStorage := postgres.NewUserStorage(
		postgresStorage.GetPool(),
	)

	tokenManager, err := jwt.NewJWT(
		secret,
		cfg.TokensTTl,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"initialize token manager: %w",
			err,
		)
	}

	limiterManager, err := redis.NewLimiter(
		cfg.Limiter,
		redisClient,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"initialize rate limiter: %w",
			err,
		)
	}

	passwordHasher := hash.Bcrypt{}
	tokenHasher := hash.SHA256{}

	userService := applicationCore.NewUserService(
		userStorage,
		application.redis,
		passwordHasher,
		limiterManager,
	)

	authService := applicationCore.NewAuthService(
		tokenManager,
		application.redis,
		application.redis,
		limiterManager,
		userStorage,
		passwordHasher,
		tokenHasher,
	)

	tokenParser, err := sharedjwt.NewParser(secret)
	if err != nil {
		return nil, fmt.Errorf(
			"initialize token parser: %w",
			err,
		)
	}

	var sessionValidator sharedsession.Validator = application.redis

	application.grpcApp = infrastructure.NewGRPC(
		logger.Named("grpc"),
		userService,
		authService,
		cfg.GRPC.Port,
		tokenParser,
		sessionValidator,
	)

	initialized = true

	return application, nil
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
		return fmt.Errorf("start application: %w", err)
	}

	return nil
}

func (a *App) Stop(ctx context.Context) error {
	a.logger.Info("application stopping")

	var grpcErr error

	if a.grpcApp != nil {
		grpcErr = a.grpcApp.Stop(ctx)
	}

	cleanupErr := a.closeResources()

	a.logger.Info("application stopped")

	return errors.Join(grpcErr, cleanupErr)
}

func (a *App) closeResources() error {
	var redisErr error

	if a.redis != nil {
		redisErr = a.redis.Stop()

		if redisErr != nil {
			a.logger.Error(
				"failed to close redis connection",
				zap.Error(redisErr),
			)
		}
	}

	if a.postgres != nil {
		a.postgres.ClosePool()
	}

	return redisErr
}
