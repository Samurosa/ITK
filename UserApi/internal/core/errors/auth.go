package errors

import "errors"

var (
	ErrRefreshExpired          = errors.New("refresh expired")
	ErrGenerateToken           = errors.New("generate access token")
	ErrInvalidToken            = errors.New("invalid access token")
	ErrGenerateTokenProcessing = errors.New("generate token processing")

	ErrSessionExpired  = errors.New("session expired")
	ErrSessionNotFound = errors.New("session not found")

	ErrInvalidContext = errors.New("invalid context")

	ErrInvalidLoginCredentials = errors.New("invalid login credentials")
	Unauthorized               = errors.New("unauthorized")
	ErrNoAccess                = errors.New("no access")
	ErrTooManyRequests         = errors.New("too many requests")

	ErrSyncRedis = errors.New("sync redis error")
)
