package jwt

import (
	"ITK_Code/m/v2/internal/core/auth"
	"ITK_Code/m/v2/internal/core/user"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type Token struct {
	log       *zap.Logger
	jwtConfig auth.JWTConfig
}

func NewJWT(log *zap.Logger, jwtConfig auth.JWTConfig) *Token {
	return &Token{
		log:       log,
		jwtConfig: jwtConfig,
	}
}

func (j *Token) Generate(user user.User, deviceID string) (auth.TokensModel, auth.AccessTokenParse, auth.RefreshTokenParse, error) {
	accessTokenString, accessToken, err := generateAccessToken(
		j.jwtConfig.Secret,
		j.jwtConfig.AccessTokenTTL,
		user,
		deviceID,
	)
	if err != nil {
		return auth.TokensModel{}, auth.AccessTokenParse{}, auth.RefreshTokenParse{}, err
	}

	refreshTokenString, refreshToken, err := generateRefreshToken(
		j.jwtConfig.Secret,
		j.jwtConfig.AccessTokenTTL,
		accessToken.Jti,
	)
	if err != nil {
		return auth.TokensModel{}, auth.AccessTokenParse{}, auth.RefreshTokenParse{}, err
	}

	return auth.TokensModel{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,

		AccessExpiresAt:  time.Now().Add(j.jwtConfig.AccessTokenTTL),
		AccessCreatedAt:  time.Now(),
		RefreshExpiresAt: time.Now().Add(j.jwtConfig.AccessTokenTTL),
		RefreshCreatedAt: time.Now(),
		RefreshTTL:       j.jwtConfig.RefreshTokenTTL,
	}, accessToken, refreshToken, nil
}

func (j *Token) ParseAccessToken(accessToken string) (auth.AccessTokenParse, error) {
	log := j.log.Named("Parse Access Token")
	token, err := jwt.ParseWithClaims(
		accessToken,
		&auth.AccessTokenParse{},
		func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, auth.ErrInvalidToken
			}
			return []byte(j.jwtConfig.Secret), nil
		},
	)
	if err != nil {
		log.Error("Parse Access Token Error", zap.Error(err))
		return auth.AccessTokenParse{}, auth.ErrInvalidToken
	}

	claims, err := GetClaimsWithAccessToken(log, token)
	if err != nil {
		return auth.AccessTokenParse{}, err
	}

	return *claims, nil
}

func (j *Token) ParseRefreshToken(refreshToken string) (auth.RefreshTokenParse, error) {
	log := j.log.Named("Parse Refresh Token")
	token, err := jwt.ParseWithClaims(
		refreshToken,
		&auth.RefreshTokenParse{},
		func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, auth.ErrInvalidToken
			}
			return []byte(j.jwtConfig.Secret), nil
		},
	)
	if err != nil {
		log.Error("Parse Refresh Token Error", zap.Error(err))
		return auth.RefreshTokenParse{}, auth.ErrInvalidToken
	}

	claims, err := GetClaimsWithRefreshToken(log, token)

	if err != nil {
		return auth.RefreshTokenParse{}, err
	}

	return *claims, nil
}
