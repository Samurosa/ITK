package jwt

import (
	"ITK_Code/m/v2/internal/core/corerrors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

func (c RefreshToken) Validate() error {
	if _, err := uuid.Parse(c.AccessTokenJTI); err != nil {
		return corerrors.ErrInvalidToken
	}
	if _, err := uuid.Parse(c.RefreshTokenJTI); err != nil {
		return corerrors.ErrInvalidToken
	}
	return nil
}
