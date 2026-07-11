package auth

import (
	"ITK_Code/m/v2/internal/core/user"
	"context"
)

type SessionRepository interface {
	Create(ctx context.Context, SessionModel SessionModel) error

	GetByUserAndDevice(ctx context.Context, userID string, deviceID string) (SessionModel, error)

	GetAllByUser(ctx context.Context, userID string) ([]SessionModel, error)

	GetByRefreshToken(ctx context.Context, refreshTokenHash []byte) (SessionModel, error)

	Update(ctx context.Context, SessionModel SessionModel) error

	DeleteByUserAndDevice(ctx context.Context, userID, deviceID string) error

	DeleteByUser(ctx context.Context, userID string) error

	DeleteExpiredSessions(ctx context.Context) error
}

type TokenManager interface {
	Generate(user user.User, deviceID string) (TokensModel, error)

	ParseAccessToken(accessToken string) (TokenParse, error)
}
