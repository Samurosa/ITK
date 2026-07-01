package auth

import "time"

type JWTConfig struct {
	Secret string

	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type SessionModel struct {
	UserID string

	DeviceID string

	RefreshTokenHash []byte

	ExpiresAt time.Time
}

type TokensModel struct {
	AccessToken  string
	RefreshToken string

	AccessExpiresAt  time.Time
	RefreshExpiresAt time.Time
}
