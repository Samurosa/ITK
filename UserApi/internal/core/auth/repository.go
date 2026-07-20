package auth

import (
	"ITK_Code/m/v2/internal/core/user"
	"context"
)

type SessionRepository interface {
	Create(ctx context.Context, jti string, SessionModel SessionModel) error

	GetByJTI(ctx context.Context, jti string) (SessionModel, error)

	Update(ctx context.Context, storedJTI string, jti string, SessionModel SessionModel) error

	DeleteByJTI(ctx context.Context, jti string, userID string) error

	DeleteByUser(ctx context.Context, userID string) error
}

type SyncPrimitiveForRedis interface {
	AcquireRefreshLock(ctx context.Context, userID string) (bool, error)
	ReleaseRefreshLock(ctx context.Context, userID string) error
}

type TokenManager interface {
	Generate(user user.User, deviceID string) (TokensModel, AccessTokenParse, RefreshTokenParse, error)

	ParseAccessToken(token string) (AccessTokenParse, error)
	ParseRefreshToken(refreshToken string) (RefreshTokenParse, error)
}
