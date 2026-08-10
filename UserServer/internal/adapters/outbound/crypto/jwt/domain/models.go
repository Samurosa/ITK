package domaintype

import (
	"github.com/golang-jwt/jwt/v5"
)

type AccessTokenParse struct {
	UserID string
	Role   string
	Device string
	Jti    string

	jwt.RegisteredClaims
}

type RefreshTokenParse struct {
	AccessTokenJTI  string
	RefreshTokenJTI string

	jwt.RegisteredClaims
}
