package app

import (
	"ITK_Code/m/v2/internal/adapters/outbound/crypto/jwt"
	"ITK_Code/m/v2/internal/adapters/outbound/repository/postgres"
	"ITK_Code/m/v2/internal/adapters/outbound/repository/redis"
	"ITK_Code/m/v2/internal/application"
	"ITK_Code/m/v2/internal/config"
	"ITK_Code/m/v2/internal/core/auth"
	"ITK_Code/m/v2/internal/core/dto"
	"ITK_Code/m/v2/internal/core/user"
	"ITK_Code/m/v2/internal/core/wallet"

	"go.uber.org/zap"
)

type Dependencies struct {
	tokenManager   auth.TokenManager
	userService    user.Service
	authService    auth.Service
	walletService  wallet.Service
	sessionStorage auth.SessionRepository
}

func NewServices(log *zap.Logger,
	cfg *config.Config,
	postgresStorage *postgres.Storage,
	redisStorage *redis.Storage,
	secret string,
) *Dependencies {

	walletStorage := postgres.NewBalanceStorage(postgresStorage.GetPool())
	userStorage := postgres.NewUserStorage(postgresStorage.GetPool())

	tokenManager := jwt.NewJWT(log,
		dto.JWTConfig{
			Secret:          secret,
			AccessTokenTTL:  cfg.TokenTTl.AccessTokenTTL,
			RefreshTokenTTL: cfg.TokenTTl.RefreshTokenTTL,
		})

	rateLimiter := redis.NewLimiter(log, cfg.Limiter, redisStorage.GetClient())

	userService := application.NewUserService(log,
		userStorage,
		redisStorage,
	)

	authService := application.NewAuthService(log,
		tokenManager,
		redisStorage,
		redisStorage,
		rateLimiter,
		userStorage,
	)

	walletService := application.NewWalletService(log,
		walletStorage,
		userStorage,
	)

	return &Dependencies{
		tokenManager:   tokenManager,
		userService:    userService,
		authService:    authService,
		walletService:  walletService,
		sessionStorage: redisStorage,
	}
}
func (a *Dependencies) TokenManager() auth.TokenManager {
	return a.tokenManager
}
func (a *Dependencies) UserService() user.Service {
	return a.userService
}
func (a *Dependencies) AuthService() auth.Service {
	return a.authService
}
func (a *Dependencies) WalletService() wallet.Service {
	return a.walletService
}
func (a *Dependencies) SessionStorage() auth.SessionRepository {
	return a.sessionStorage
}
