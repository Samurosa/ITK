package auth

import (
	"time"
)

type SessionModel struct {
	UserID   string
	DeviceID string

	RefreshTokenHash string
	TTL              time.Duration
	ExpiresAt        time.Time
	CreatedAt        time.Time
}
