package jwt

import (
	"ITK_Code/m/v2/internal/core/auth"
	"ITK_Code/m/v2/internal/core/user"
	"time"
)

type Token struct {
	jwtConfig auth.JWTConfig
}

func NewJWT(jwtConfig auth.JWTConfig) *Token {
	return &Token{jwtConfig: jwtConfig}
}

func (j *Token) Generate(user user.User) (auth.TokensModel, error) {
	refreshTokenString, err := generateRefreshToken()
	if err != nil {
		return auth.TokensModel{}, err
	}

	accessTokenString, err := generateAccessToken(j.jwtConfig.Secret, j.jwtConfig.AccessTokenTTL, user)
	if err != nil {
		return auth.TokensModel{}, err
	}

	return auth.TokensModel{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,

		AccessExpiresAt:  time.Now().Add(j.jwtConfig.AccessTokenTTL),
		RefreshExpiresAt: time.Now().Add(j.jwtConfig.RefreshTokenTTL),
	}, nil
}
