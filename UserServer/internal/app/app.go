package app

import (
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
	logger, err := zap.NewProduction()
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

	walletStorage := postgres.NewBalanceStorage(postgresStorage.GetPool())

	tokenManager, err := jwt.NewJWT(secret, cfg.TokensTTl)
	if err != nil {
		postgresStorage.ClosePool()
		cancel()
		return nil, err
	}

	limiterManager := redis.NewLimiter(log, cfg.Limiter, redisClient)

	user := application.NewUserService(log, userStorage, redisStorage)
	auth := application.NewAuthService(log, tokenManager, redisStorage, redisStorage, limiterManager, userStorage)
	wallet := application.NewWalletService(log, walletStorage, userStorage)

	tokenParser, err := sharedjwt.NewParser(secret)
	if err != nil {
		redisStorage.Stop()
		postgresStorage.ClosePool()
		cancel()
		return nil, err
	}

	var sessionValidator sharedsession.Validator = redisStorage

	grpcApp := infrastructure.NewGRPC(
		logger,
		user,
		auth,
		wallet,
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
