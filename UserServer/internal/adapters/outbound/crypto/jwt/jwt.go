package jwt

import (
	"ITK_Code/m/v2/internal/config"
	"ITK_Code/m/v2/internal/core/corerrors"
	"ITK_Code/m/v2/internal/core/dto"
	"ITK_Code/m/v2/internal/core/user"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Token struct {
	secret    string
	tokensTTL config.TokensTTL
}

func NewJWT(secret string, tokensTTl config.TokensTTL) (*Token, error) {

	if err := configValidate(secret, tokensTTl); err != nil {
		return nil, err
	}

	return &Token{
		secret:    secret,
		tokensTTL: tokensTTl,
	}, nil
}

func (j *Token) Generate(user user.User, deviceID string) (dto.TokensModel, dto.AccessTokenParse, dto.RefreshTokenParse, error) {
	accessTokenString, accessToken, err := generateAccessToken(
		j.secret,
		j.tokensTTL.AccessTokenTTL,
		user,
		deviceID,
	)
	if err != nil {
		return dto.TokensModel{}, dto.AccessTokenParse{}, dto.RefreshTokenParse{}, err
	}

	refreshTokenString, refreshToken, err := generateRefreshToken(
		j.secret,
		j.tokensTTL.RefreshTokenTTL,
		accessToken.Jti,
	)
	if err != nil {
		return dto.TokensModel{}, dto.AccessTokenParse{}, dto.RefreshTokenParse{}, err
	}

	return dto.TokensModel{
			AccessToken:  accessTokenString,
			RefreshToken: refreshTokenString,

			AccessExpiresAt:  time.Now().Add(j.tokensTTL.AccessTokenTTL),
			AccessIssuedAt:   time.Now(),
			RefreshExpiresAt: time.Now().Add(j.tokensTTL.RefreshTokenTTL),
			RefreshIssuedAt:  time.Now(),
			RefreshTTL:       j.tokensTTL.RefreshTokenTTL,
		},
		dto.AccessTokenParse{
			UserID: accessToken.UserID,
			Role:   accessToken.Role,
			Device: deviceID,
			Jti:    accessToken.Jti,
		},
		dto.RefreshTokenParse{
			AccessTokenJTI:  refreshToken.AccessTokenJTI,
			RefreshTokenJTI: refreshToken.RefreshTokenJTI,
		}, nil
}

func (j *Token) ParseRefreshToken(refreshToken string) (dto.RefreshTokenParse, error) {
	token, err := jwt.ParseWithClaims(
		refreshToken,
		&RefreshTokenParse{},
		func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, corerrors.ErrInvalidToken
			}
			return []byte(j.secret), nil
		},
	)
	if err != nil {
		return dto.RefreshTokenParse{}, corerrors.ErrInvalidToken
	}

	claims, err := GetClaimsWithRefreshToken(token)

	if err != nil {
		return dto.RefreshTokenParse{}, err
	}

	return *claims, nil
}

func configValidate(secret string, tokensTTl config.TokensTTL) error {
	if secret == "" {
		return corerrors.ErrJWTSecret
	}
	if tokensTTl.AccessTokenTTL < 1 {
		return corerrors.ErrAccessTokenTTL
	}
	if tokensTTl.RefreshTokenTTL < 1 {
		return corerrors.ErrRefreshTokenTTL
	}
	return nil
}
