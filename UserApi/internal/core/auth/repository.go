package auth

import (
	"ITK_Code/m/v2/internal/core/user"
	"context"
)

type SessionRepository interface {
	Create(ctx context.Context, jti string, SessionModel SessionModel) error

	GetByJTI(ctx context.Context, jti string) (SessionModel, error)

	GetAllByUser(ctx context.Context, userID string) ([]SessionModel, error)

	GetByRefreshToken(ctx context.Context, jti string) (SessionModel, error)

	Update(ctx context.Context, SessionModel SessionModel) error

	DeleteByJTI(ctx context.Context, jti string) error

	DeleteByUser(ctx context.Context, userID string) error

	DeleteExpiredSessions(ctx context.Context) error
}

type TokenManager interface {
	Generate(user user.User, deviceID string) (TokensModel, error)

	ParseAccessToken(token string) (AccessTokenParse, error)
	ParseRefreshToken(refreshToken string) (RefreshTokenParse, error)
}
