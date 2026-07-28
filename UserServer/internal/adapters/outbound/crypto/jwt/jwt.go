package jwt

import (
	"ITK_Code/m/v2/internal/core/dto"
	"ITK_Code/m/v2/internal/core/errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type Token struct {
	log       *zap.Logger
	jwtConfig dto.JWTConfig
}

func NewJWT(log *zap.Logger, jwtConfig dto.JWTConfig) *Token {
	return &Token{
		log:       log,
		jwtConfig: jwtConfig,
	}
}

func (j *Token) Generate(user dto.User, deviceID string) (dto.TokensModel, dto.AccessTokenParse, dto.RefreshTokenParse, error) {
	accessTokenString, accessToken, err := generateAccessToken(
		j.jwtConfig.Secret,
		j.jwtConfig.AccessTokenTTL,
		user,
		deviceID,
	)
	if err != nil {
		return dto.TokensModel{}, dto.AccessTokenParse{}, dto.RefreshTokenParse{}, err
	}

	refreshTokenString, refreshToken, err := generateRefreshToken(
		j.jwtConfig.Secret,
		j.jwtConfig.RefreshTokenTTL,
		accessToken.Jti,
	)
	if err != nil {
		return dto.TokensModel{}, dto.AccessTokenParse{}, dto.RefreshTokenParse{}, err
	}

	return dto.TokensModel{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,

		AccessExpiresAt:  time.Now().Add(j.jwtConfig.AccessTokenTTL),
		AccessCreatedAt:  time.Now(),
		RefreshExpiresAt: time.Now().Add(j.jwtConfig.RefreshTokenTTL),
		RefreshCreatedAt: time.Now(),
		RefreshTTL:       j.jwtConfig.RefreshTokenTTL,
	}, accessToken, refreshToken, nil
}

func (j *Token) ParseAccessToken(accessToken string) (dto.AccessTokenParse, error) {
	log := j.log.Named("Parse Access Token")
	token, err := jwt.ParseWithClaims(
		accessToken,
		&dto.AccessTokenParse{},
		func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, errors.ErrInvalidToken
			}
			return []byte(j.jwtConfig.Secret), nil
		},
	)
	if err != nil {
		log.Error("Parse Access Token Error", zap.Error(err))
		return dto.AccessTokenParse{}, errors.ErrInvalidToken
	}

	claims, err := GetClaimsWithAccessToken(log, token)
	if err != nil {
		return dto.AccessTokenParse{}, err
	}

	return *claims, nil
}

func (j *Token) ParseRefreshToken(refreshToken string) (dto.RefreshTokenParse, error) {
	log := j.log.Named("Parse Refresh Token")
	token, err := jwt.ParseWithClaims(
		refreshToken,
		&dto.RefreshTokenParse{},
		func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, errors.ErrInvalidToken
			}
			return []byte(j.jwtConfig.Secret), nil
		},
	)
	if err != nil {
		log.Error("Parse Refresh Token Error", zap.Error(err))
		return dto.RefreshTokenParse{}, errors.ErrInvalidToken
	}

	claims, err := GetClaimsWithRefreshToken(log, token)

	if err != nil {
		return dto.RefreshTokenParse{}, err
	}

	return *claims, nil
}
