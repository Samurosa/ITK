package ports

import (
	"ITK_Code/m/v2/internal/core/auth"
	"ITK_Code/m/v2/internal/core/dto"
	"ITK_Code/m/v2/internal/core/user"
	"context"
)

type SessionRepository interface {
	Create(context.Context, string, auth.SessionModel) error
	GetByJTI(context.Context, string) (auth.SessionModel, error)
	Update(context.Context, string, string, auth.SessionModel) error
	DeleteByJTI(context.Context, string, string) error
	DeleteByUser(context.Context, string) error
}

type TokenManager interface {
	Generate(user.User, string) (dto.TokensModel, dto.AccessToken, dto.RefreshToken, error)
	ParseRefreshToken(string) (dto.RefreshToken, error)
}

type RefreshLock interface {
	AcquireRefreshLock(context.Context, string) (string, error)
	ReleaseRefreshLock(context.Context, string, string) error
}

type RateLimiter interface {
	Allow(context.Context, string, string) (bool, error)
}

type PasswordChangeLimiter interface {
	AllowPasswordChange(context.Context, string) (bool, error)
}

type TokenHasher interface {
	GenerateHashSHA256(string) string
	CompareHashSHA256(string, string) error
}
