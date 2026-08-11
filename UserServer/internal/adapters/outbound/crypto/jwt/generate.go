package jwt

import (
	"ITK_Code/m/v2/internal/core/dto"
	"ITK_Code/m/v2/internal/core/errors"
	"ITK_Code/m/v2/internal/core/user"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func generateRefreshToken(
	secret string,
	refreshTokenTTL time.Duration,
	jti string,
) (string, RefreshTokenParse, error) {

	claimsRefreshToken := RefreshTokenParse{
		AccessTokenJTI:  jti,
		RefreshTokenJTI: uuid.NewString(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenTTL)),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsRefreshToken)

	refreshTokenString, err := refreshToken.SignedString([]byte(secret))
	if err != nil {
		return "", RefreshTokenParse{}, errors.ErrGenerateToken
	}

	return refreshTokenString, claimsRefreshToken, nil
}

func generateAccessToken(
	secret string,
	accessTokenTTL time.Duration,
	user user.User,
	deviceId string,
) (string, AccessTokenParse, error) {

	tokenID := uuid.NewString()

	claimsAccessToken := AccessTokenParse{
		UserID: user.ID,
		Role:   string(user.Role),
		Device: deviceId,
		Jti:    tokenID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenTTL)),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsAccessToken)

	accessTokenString, err := accessToken.SignedString([]byte(secret))
	if err != nil {
		return "", AccessTokenParse{}, errors.ErrGenerateToken
	}

	return accessTokenString, claimsAccessToken, nil
}

func GetClaimsWithAccessToken(log *zap.Logger, token *jwt.Token) (*dto.AccessTokenParse, error) {

	if !token.Valid {
		log.Error("token is not valid")
		return &dto.AccessTokenParse{}, errors.ErrInvalidToken
	}

	claims, ok := token.Claims.(*AccessTokenParse)
	if !ok {
		log.Error("token claims is not found")
		return &dto.AccessTokenParse{}, errors.ErrInvalidToken
	}

	return &dto.AccessTokenParse{
		UserID: claims.UserID,
		Role:   claims.Role,
		Device: claims.Device,
		Jti:    claims.Jti,
	}, nil
}

func GetClaimsWithRefreshToken(log *zap.Logger, token *jwt.Token) (*dto.RefreshTokenParse, error) {

	if !token.Valid {
		log.Error("token is not valid")
		return &dto.RefreshTokenParse{}, errors.ErrInvalidToken
	}

	claims, ok := token.Claims.(*RefreshTokenParse)
	if !ok {
		log.Error("token claims is not found")
		return &dto.RefreshTokenParse{}, errors.ErrInvalidToken
	}

	return &dto.RefreshTokenParse{
		AccessTokenJTI:  claims.AccessTokenJTI,
		RefreshTokenJTI: claims.RefreshTokenJTI,
	}, nil
}
