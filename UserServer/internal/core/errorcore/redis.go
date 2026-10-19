package errorcore

import (
	"errors"
)

var (
	ErrSyncRedis   = errors.New("sync redis error")
	ErrPingToRedis = errors.New("ping to redis error")
)
