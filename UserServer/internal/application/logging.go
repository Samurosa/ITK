package application

import (
	"ITK_Code/m/v2/internal/core/auth"
	"ITK_Code/m/v2/internal/core/user"
	"context"
	"errors"

	"go.uber.org/zap"
)

// Expected lookup misses and canceled requests are not infrastructure failures.
func logOperationError(log *zap.Logger, message string, err error) {
	switch {
	case errors.Is(err, user.ErrUserNotFound), errors.Is(err, user.ErrEmailIsExist),
		errors.Is(err, auth.ErrSessionNotFound), errors.Is(err, context.Canceled):
		log.Debug(message, zap.Error(err))
	case errors.Is(err, context.DeadlineExceeded):
		log.Warn(message, zap.Error(err))
	default:
		log.Error(message, zap.Error(err))
	}
}

func logPasswordVerificationError(log *zap.Logger, err error) {
	if errors.Is(err, auth.ErrIncorrectPassword) {
		log.Warn("password verification rejected")
		return
	}
	logOperationError(log, "password verification failed", err)
}
