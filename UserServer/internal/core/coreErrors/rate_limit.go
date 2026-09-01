package coreErrors

import "errors"

var (
	ErrTooManyRequests = errors.New("too many requests")
)
