package dto

import (
	"time"
)

type JWTConfig struct {
	Secret string

	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type TokensModel struct {
	AccessToken  string
	RefreshToken string

	AccessExpiresAt  time.Time
	AccessIssuedAt   time.Time
	RefreshExpiresAt time.Time
	RefreshIssuedAt  time.Time
	RefreshTTL       time.Duration
}

type AccessTokenParse struct {
	UserID string
	Role   string
	Device string
	Jti    string
}

type RefreshTokenParse struct {
	AccessTokenJTI  string
	RefreshTokenJTI string
}
