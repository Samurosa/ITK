package auth

import "errors"

var (
	ErrRefreshExpired       = errors.New("refresh expired")
	ErrGenerateRefreshToken = errors.New("generate refresh token")
	ErrGenerateAccessToken  = errors.New("generate access token")

	Unauthorized = errors.New("incorrect login or password")
	ErrNoAccess  = errors.New("no access")
)
