package application

import (
	authCore "ITK_Code/m/v2/internal/core/auth"
	userCore "ITK_Code/m/v2/internal/core/user"
	"ITK_Code/m/v2/internal/core/wallet"

	"go.uber.org/zap"
)

type User struct {
	log *zap.Logger

	tokenManager authCore.TokenManager

	sessionStorage authCore.SessionRepository

	userSaver    userCore.Save
	userProvider userCore.Provider

	balanceRepository wallet.Repository
}

func New(
	log *zap.Logger,
	tokenManager authCore.TokenManager,
	sessionStorage authCore.SessionRepository,
	userSaver userCore.Save,
	userProvider userCore.Provider,
	balanceRepository wallet.Repository,
) *User {
	return &User{
		log:               log,
		tokenManager:      tokenManager,
		sessionStorage:    sessionStorage,
		userSaver:         userSaver,
		userProvider:      userProvider,
		balanceRepository: balanceRepository,
	}
}
