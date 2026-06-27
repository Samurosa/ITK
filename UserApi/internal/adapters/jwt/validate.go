package jwt

import (
	"ITK_Code/m/v2/internal/core/auth"
	"time"
)

func validateExpiration(session auth.SessionModel) error {
	if session.ExpiresAt.Before(time.Now()) {
		return auth.ErrRefreshExpired
	}
	return nil
}
