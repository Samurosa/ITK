package errors

import "errors"

var (
	ErrRefreshExpired          = errors.New("refresh expired")
	ErrGenerateToken           = errors.New("generate access token")
	ErrInvalidToken            = errors.New("invalid access token")
	ErrGenerateTokenProcessing = errors.New("generate token processing")
)
