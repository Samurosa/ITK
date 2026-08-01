package dto

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTConfig struct {
	Secret string

	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type SessionModel struct {
	UserID   string
	DeviceID string

	RefreshTokenHash string
	TTL              time.Duration
	ExpiresAt        time.Time
	CreatedAt        time.Time
}

type TokensModel struct {
	AccessToken  string
	RefreshToken string

	AccessExpiresAt  time.Time
	AccessCreatedAt  time.Time
	RefreshExpiresAt time.Time
	RefreshCreatedAt time.Time
	RefreshTTL       time.Duration
}

type AccessTokenParse struct {
	UserID string
	Role   string
	Device string
	Jti    string

	jwt.RegisteredClaims
}

type RefreshTokenParse struct {
	AccessTokenJTI  string
	RefreshTokenJTI string

	jwt.RegisteredClaims
}
