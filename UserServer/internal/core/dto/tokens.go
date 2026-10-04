package dto

import (
	"time"
)

type TokensModel struct {
	AccessToken  string
	RefreshToken string

	AccessExpiresAt  time.Time
	AccessIssuedAt   time.Time
	RefreshExpiresAt time.Time
	RefreshIssuedAt  time.Time
	RefreshTTL       time.Duration
}

type AccessToken struct {
	UserID   string `json:"user_id"`
	Role     string `json:"role"`
	DeviceID string `json:"device_id"`
	JTI      string `json:"jti"`
}

type RefreshToken struct {
	AccessTokenJTI  string
	RefreshTokenJTI string
}
