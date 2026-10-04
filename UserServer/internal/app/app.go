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
		log.Error("failed to initialize token manager", zap.Error(err))
		if closeErr := redisStorage.Stop(); closeErr != nil {
			log.Warn("Failed to stop redis storage", zap.Error(closeErr))
		}
		postgresStorage.ClosePool()
		cancel()
		return nil, err
	}

	limiterManager, err := redis.NewLimiter(cfg.Limiter, redisClient)
	if err != nil {
		log.Error("failed to initialize rate limiter", zap.Error(err))
		if closeErr := redisStorage.Stop(); closeErr != nil {
			log.Warn("Failed to stop redis storage", zap.Error(closeErr))
		}
		postgresStorage.ClosePool()
		cancel()
		return nil, err
	}

	passwordHasher := hash.Bcrypt{}
	tokenHasher := hash.SHA256{}
	user := application.NewUserService(userStorage, redisStorage, passwordHasher, limiterManager)
	auth := application.NewAuthService(tokenManager, redisStorage, redisStorage, limiterManager, userStorage, passwordHasher, tokenHasher)

	tokenParser, err := sharedjwt.NewParser(secret)
	if err != nil {
		log.Error("failed to initialize token parser", zap.Error(err))
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
				"grpc server failed",
				zap.Error(err),
			)
		}
	}()
}

func (app *App) Stop() {
	app.logger.Info("application stopping")

	app.grpcApp.Stop()
	app.cancel()

	app.postgres.ClosePool()

	if err := app.redis.Stop(); err != nil {
		app.logger.Error(
			"failed to close redis connection",
			zap.Error(err),
		)
	}

	app.logger.Info("application stopped")

	if err := app.logger.Sync(); err != nil {
		app.logger.Error(
			"failed to sync logger",
			zap.Error(err),
		)
	}
}
