package auth

import "errors"

var (
	ErrSessionExpired  = errors.New("session expired")
	ErrSessionNotFound = errors.New("session not found")

	ErrIncorrectCredentials = errors.New("incorrect login or password")
	ErrUnauthorized         = errors.New("unauthorized")
	ErrNoAccess             = errors.New("no access")
)
