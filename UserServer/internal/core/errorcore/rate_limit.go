package errorcore

import "errors"

var (
	ErrInvalidLimiterConfig = errors.New("limiter capacity must be positive and timer at least 1ms")
	ErrTooManyRequests      = errors.New("too many requests")
)
