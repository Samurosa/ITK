package jwt

import (
	"ITK_Code/m/v2/internal/core/auth"
	"ITK_Code/m/v2/internal/core/user"
	"ITK_Code/m/v2/internal/core/user/models"
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type JWT struct {
	jwtConfig auth.JWTConfig
}

func NewJWT(jwtConfig auth.JWTConfig) *JWT {
	return &JWT{jwtConfig: jwtConfig}
}

func (j *JWT) GenerateTokens(ctx context.Context, user models.User, deviceID string) (auth.TokensModel, error) {
	refreshTokenString, err := generateRefreshToken()
	if err != nil {
		return auth.TokensModel{}, err
	}

	refreshTokenHash, err := bcrypt.GenerateFromPassword([]byte(refreshTokenString), bcrypt.DefaultCost)
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

func (j *JWT) RefreshAccessToken(ctx context.Context) (auth.TokensModel, error) {

	return auth.TokensModel{}, nil
}
