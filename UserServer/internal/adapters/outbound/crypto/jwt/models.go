package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
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

	jwt.RegisteredClaims
}

type RefreshToken struct {
	AccessTokenJTI  string
	RefreshTokenJTI string

	jwt.RegisteredClaims
}
