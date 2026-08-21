package auth

import "errors"

var (
	ErrSessionNotFound = errors.New("session not found")

	ErrIncorrectCredentials = errors.New("incorrect login or password")
	ErrIncorrectPassword    = errors.New("incorrect password")
	ErrUnauthorized         = errors.New("unauthorized")
	ErrNoAccess             = errors.New("no access")
)
