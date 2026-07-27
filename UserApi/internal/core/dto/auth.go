package dto

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type ContextKey string

const (
	UserIDContextKey ContextKey = "user_id"
	RoleContextKey   ContextKey = "role"
	DeviceContextKey ContextKey = "device"
	JtiContextKey    ContextKey = "jti"
	ClientIPKey      ContextKey = "client_ip"
)

type JWTConfig struct {
	Secret string

	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type SessionModel struct {
	UserID   string `redis:"user_id"`
	DeviceID string `redis:"device_id"`

	RefreshTokenHash string `redis:"refresh_token_hash"`

	TTL       time.Duration `redis:"ttl"`
	ExpiresAt time.Time     `redis:"expires_at"`
	CreatedAt time.Time     `redis:"created_at"`
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
