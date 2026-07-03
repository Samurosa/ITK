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

func (j *Token) ParseAccessToken(accessToken string) (auth.TokenParse, error) {
	log := j.log.Named("Parse Access Token")
	token, err := jwt.ParseWithClaims(
		accessToken,
		&auth.TokenParse{},
		func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, auth.ErrInvalidToken
			}
			return []byte(j.jwtConfig.Secret), nil
		},
	)
	if err != nil {
		log.Error("Parse Access Token Error", zap.Error(err))
		return auth.TokenParse{}, err
	}

	claims, err := GetClaimsWithToken(log, token)

	if err != nil {
		return auth.TokenParse{}, err
	}

	return *claims, nil
}
