package jwt

import (
	"ITK_Code/m/v2/internal/core/auth"
	"ITK_Code/m/v2/internal/core/user/models"
	"time"
)

type JWT struct {
	jwtConfig auth.JWTConfig
}

func NewJWT(jwtConfig auth.JWTConfig) *JWT {
	return &JWT{jwtConfig: jwtConfig}
}

func (j *JWT) GenerateTokens(user models.User) (auth.TokensModel, error) {
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
