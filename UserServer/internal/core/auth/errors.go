package auth

import (
	"errors"

	"github.com/Samurosa/exchange-common/shared/auth/session"
)

var (
	ErrSessionNotFound = session.ErrNotFound

	ErrIncorrectCredentials = errors.New("incorrect login or password")
	ErrIncorrectPassword    = errors.New("incorrect password")
	ErrNoAccess             = errors.New("no access")
)
