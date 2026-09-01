package coreErrors

import (
	"errors"
)

var (
	ErrSyncRedis = errors.New("sync redis error")
)
