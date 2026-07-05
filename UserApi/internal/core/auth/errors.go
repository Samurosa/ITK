package auth

import "errors"

var (
	ErrRefreshExpired = errors.New("refresh expired")
	ErrGenerateToken  = errors.New("generate access token")
	ErrInvalidToken   = errors.New("access token parse error")

	Unauthorized = errors.New("incorrect login or password")
	ErrNoAccess  = errors.New("no access")
)
