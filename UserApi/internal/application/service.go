package application

import (
	authCore "ITK_Code/m/v2/internal/core/auth"
	userCore "ITK_Code/m/v2/internal/core/user"

	"go.uber.org/zap"
)

type User struct {
	log *zap.Logger

	tokenManager authCore.TokenManager

	sessionStorage authCore.SessionRepository

	userSaver    userCore.Save
	userProvider userCore.Provider
}

func New(
	log *zap.Logger,
	tokenManager authCore.TokenManager,
	sessionStorage authCore.SessionRepository,
	userSaver userCore.Save,
	userProvider userCore.Provider,
) *User {
	return &User{
		log:            log,
		tokenManager:   tokenManager,
		sessionStorage: sessionStorage,
		userSaver:      userSaver,
		userProvider:   userProvider,
	}
}
