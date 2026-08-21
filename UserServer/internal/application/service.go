package application

import (
	authCore "ITK_Code/m/v2/internal/core/auth"
	userCore "ITK_Code/m/v2/internal/core/user"
	"ITK_Code/m/v2/internal/core/wallet"

	"go.uber.org/zap"
)

type User struct {
	log *zap.Logger

	userRepository userCore.Repository

	sessionStorage authCore.SessionRepository
}

type Auth struct {
	log *zap.Logger

	tokenManager          authCore.TokenManager
	sessionStorage        authCore.SessionRepository
	syncPrimitiveForRedis authCore.SyncPrimitiveForRedis
	rateLimiting          authCore.RateLimiting

	userRepository userCore.Repository
}

type Wallet struct {
	log *zap.Logger

	balanceRepository wallet.Repository

	userRepository userCore.Repository
}

func NewUserService(
	log *zap.Logger,
	userRepository userCore.Repository,
	sessionStorage authCore.SessionRepository,
) *User {
	return &User{
		log:            log,
		userRepository: userRepository,
		sessionStorage: sessionStorage,
	}
}

func NewAuthService(
	log *zap.Logger,
	tokenManager authCore.TokenManager,
	sessionStorage authCore.SessionRepository,
	syncPrimitiveForRedis authCore.SyncPrimitiveForRedis,
	rateLimiting authCore.RateLimiting,
	userRepository userCore.Repository,

) *Auth {
	return &Auth{
		log:                   log,
		tokenManager:          tokenManager,
		sessionStorage:        sessionStorage,
		syncPrimitiveForRedis: syncPrimitiveForRedis,
		rateLimiting:          rateLimiting,
		userRepository:        userRepository,
	}
}

func NewWalletService(
	log *zap.Logger,
	balanceRepository wallet.Repository,
	userRepository userCore.Repository,
) *Wallet {
	return &Wallet{
		log:               log,
		balanceRepository: balanceRepository,
		userRepository:    userRepository,
	}
}
