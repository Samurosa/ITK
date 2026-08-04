package application

import (
	authCore "ITK_Code/m/v2/internal/core/auth"
	userCore "ITK_Code/m/v2/internal/core/user"
	"ITK_Code/m/v2/internal/core/wallet"

	"go.uber.org/zap"
)

type User struct {
	log *zap.Logger

	userSaver    userCore.Save
	userProvider userCore.Provider

	sessionStorage authCore.SessionRepository
}

type Auth struct {
	log *zap.Logger

	tokenManager          authCore.TokenManager
	sessionStorage        authCore.SessionRepository
	syncPrimitiveForRedis authCore.SyncPrimitiveForRedis
	rateLimiting          authCore.RateLimiting

	userSaver    userCore.Save
	userProvider userCore.Provider
}

type Wallet struct {
	log *zap.Logger

	balanceRepository wallet.Repository

	userProvider userCore.Provider
}

func NewUserService(
	log *zap.Logger,
	userSaver userCore.Save,
	userProvider userCore.Provider,
	sessionStorage authCore.SessionRepository,
) *User {
	return &User{
		log:            log,
		userSaver:      userSaver,
		userProvider:   userProvider,
		sessionStorage: sessionStorage,
	}
}

func NewAuthService(
	log *zap.Logger,
	tokenManager authCore.TokenManager,
	sessionStorage authCore.SessionRepository,
	syncPrimitiveForRedis authCore.SyncPrimitiveForRedis,
	rateLimiting authCore.RateLimiting,
	userSaver userCore.Save,
	userProvider userCore.Provider,

) *Auth {
	return &Auth{
		log:                   log,
		tokenManager:          tokenManager,
		sessionStorage:        sessionStorage,
		syncPrimitiveForRedis: syncPrimitiveForRedis,
		rateLimiting:          rateLimiting,
		userSaver:             userSaver,
		userProvider:          userProvider,
	}
}

func NewWalletService(
	log *zap.Logger,
	balanceRepository wallet.Repository,
	userProvider userCore.Provider,
) *Wallet {
	return &Wallet{
		log:               log,
		balanceRepository: balanceRepository,
		userProvider:      userProvider,
	}
}
