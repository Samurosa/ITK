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
	AccessCreatedAt  time.Time
	RefreshExpiresAt time.Time
	RefreshCreatedAt time.Time
	RefreshTTL       time.Duration
}
