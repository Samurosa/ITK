package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type ContextKey string

const (
	UserIDContextKey ContextKey = "user_id"
	RoleContextKey   ContextKey = "role"
	DeviceContextKey ContextKey = "device"
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

type AccessTokenParse struct {
	UserID string
	Role   string
	Device string
	Jti    string

	jwt.RegisteredClaims
}

type RefreshTokenParse struct {
	AccessTokenJti string
	Jti            string

	jwt.RegisteredClaims
}
