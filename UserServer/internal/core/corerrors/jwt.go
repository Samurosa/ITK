package corerrors

import "errors"

var (
	ErrGenerateToken           = errors.New("generate access token")
	ErrInvalidToken            = errors.New("invalid access token")
	ErrDeviceIDEmpty           = errors.New("device ID is empty")
	ErrGenerateTokenProcessing = errors.New("generate token processing")
)
