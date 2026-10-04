package app

import (
	"ITK_Code/m/v2/internal/adapters/outbound/crypto/hash"
	"ITK_Code/m/v2/internal/adapters/outbound/crypto/jwt"
	"ITK_Code/m/v2/internal/application"
	"context"

	"ITK_Code/m/v2/internal/adapters/outbound/repository/postgres"
	"ITK_Code/m/v2/internal/adapters/outbound/repository/redis"
	"ITK_Code/m/v2/internal/config"
	"ITK_Code/m/v2/internal/infrastructure"

	sharedjwt "github.com/Samurosa/exchange-common/shared/auth/jwt"
	sharedsession "github.com/Samurosa/exchange-common/shared/auth/session"
	"go.uber.org/zap"
)

type App struct {
	logger   *zap.Logger
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
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}

	log := logger.Named("app")

	ctx, cancel := context.WithCancel(context.Background())

	postgresStorage, err := postgres.NewStorage(
		ctx,
		log,
		cfg.Postgres,
	)
	if err != nil {
		cancel()
		return nil, err
	}

	redisClient, err := redis.NewRedisClient(
		ctx,
		log,
		cfg.Redis,
	)
	if err != nil {
		postgresStorage.ClosePool()
		cancel()
		return nil, err
	}

	redisStorage := redis.NewStorage(redisClient)

	userStorage := postgres.NewUserStorage(postgresStorage.GetPool())

	tokenManager, err := jwt.NewJWT(secret, cfg.TokensTTl)
	if err != nil {
		if closeErr := redisStorage.Stop(); closeErr != nil {
			log.Warn("Failed to stop redis storage", zap.Error(closeErr))
		}
		postgresStorage.ClosePool()
		cancel()
		return nil, err
	}

	limiterManager, err := redis.NewLimiter(log, cfg.Limiter, redisClient)
	if err != nil {
		if closeErr := redisStorage.Stop(); closeErr != nil {
			log.Warn("Failed to stop redis storage", zap.Error(closeErr))
		}
		postgresStorage.ClosePool()
		cancel()
		return nil, err
	}

	passwordHasher := hash.Bcrypt{}
	tokenHasher := hash.SHA256{}
	user := application.NewUserService(log, userStorage, redisStorage, passwordHasher, limiterManager)
	auth := application.NewAuthService(log, tokenManager, redisStorage, redisStorage, limiterManager, userStorage, passwordHasher, tokenHasher)

	tokenParser, err := sharedjwt.NewParser(secret)
	if err != nil {
		if closeErr := redisStorage.Stop(); closeErr != nil {
			log.Warn("Failed to stop redis storage", zap.Error(closeErr))
		}
		postgresStorage.ClosePool()
		cancel()
		return nil, err
	}

	var sessionValidator sharedsession.Validator = redisStorage

	grpcApp := infrastructure.NewGRPC(
		logger,
		user,
		auth,
		cfg.GRPC.Port,
		tokenParser,
		sessionValidator,
	)

	return &App{
		logger:   logger,
		ctx:      ctx,
		cancel:   cancel,
		postgres: postgresStorage,
		redis:    redisStorage,
		grpcApp:  grpcApp,
	}, nil
}

func (app *App) Start() {
	log := app.logger.Named("grpc")

	go func() {
		if err := app.grpcApp.Run(); err != nil {
			log.Error(
				"grpc server stopped",
				zap.Error(err),
			)
		}
	}()
}

func (app *App) Stop() {
	app.logger.Debug("application stop")

	app.grpcApp.Stop()
	app.cancel()

	app.postgres.ClosePool()

	if err := app.redis.Stop(); err != nil {
		app.logger.Error(
			"redis stop",
			zap.Error(err),
		)
	}

	if err := app.logger.Sync(); err != nil {
		app.logger.Error(
			"logger sync",
			zap.Error(err),
		)
	}
}
