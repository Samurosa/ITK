package corerrors

import (
	"errors"
)

var (
	ErrSyncRedis       = errors.New("sync redis error")
	ErrSessionNotFound = errors.New("session not found")
)
