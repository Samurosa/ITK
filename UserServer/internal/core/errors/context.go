package errors

import "errors"

var (
	ErrInvalidContext = errors.New("invalid context")

	ErrCtxForUpdateNotFound = errors.New("base context for deletion not found")

	Canceled         = errors.New("context canceled")
	DeadlineExceeded = errors.New("deadline exceeded")
)
