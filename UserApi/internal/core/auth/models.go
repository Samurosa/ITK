package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type ContextKey string

const (
	UserIDContextKey ContextKey = "user_id"
	RoleContextKey   ContextKey = "role"
)

type JWTConfig struct {
	Secret string

	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type SessionModel struct {
	ID       string
	UserID   string
	DeviceID string

	RefreshTokenHash []byte

	ExpiresAt time.Time
	CreatedAt time.Time
}

type TokensModel struct {
	AccessToken  string
	RefreshToken string

	AccessExpiresAt  time.Time
	RefreshExpiresAt time.Time
}

type TokenParse struct {
	UserID string
	Role   string
	Jti    string

	jwt.RegisteredClaims
}
