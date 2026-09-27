package corerrors

import "errors"

var (
	ErrJWTSecret       = errors.New("JWT secret is missing from the config or has an incorrect format")
	ErrAccessTokenTTL  = errors.New("AccessTokenTTL is missing from the config or has an incorrect format")
	ErrRefreshTokenTTL = errors.New("RefreshTTL is missing from the config or has an incorrect format")
)
