package auth

import (
	"ITK_Code/m/v2/internal/core/user/models"
	"context"
)

type SessionRepository interface {
	Create(ctx context.Context, SessionModel SessionModel) error

	GetByRefreshTokenHash(ctx context.Context, hash string) (SessionModel, error)

	Update(ctx context.Context, SessionModel SessionModel) error

	DeleteByUserAndDevice(ctx context.Context, userID, deviceID string) error

	DeleteByRefreshTokenHash(ctx context.Context, hash string) error
}

type TokenManager interface {
	GenerateTokens(user models.User) (TokensModel, error)
}
