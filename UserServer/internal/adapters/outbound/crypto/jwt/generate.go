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
) (string, dto.RefreshTokenParse, error) {

	claimsRefreshToken := dto.RefreshTokenParse{
		AccessTokenJTI:  jti,
		RefreshTokenJTI: uuid.NewString(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenTTL)),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsRefreshToken)

	refreshTokenString, err := refreshToken.SignedString([]byte(secret))
	if err != nil {
		return "", dto.RefreshTokenParse{}, errors.ErrGenerateToken
	}

	return refreshTokenString, claimsRefreshToken, nil
}

func generateAccessToken(
	secret string,
	accessTokenTTL time.Duration,
	user user.User,
	deviceId string,
) (string, dto.AccessTokenParse, error) {

	tokenID := uuid.NewString()

	claimsAccessToken := dto.AccessTokenParse{
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
		return "", dto.AccessTokenParse{}, errors.ErrGenerateToken
	}

	return accessTokenString, claimsAccessToken, nil
}

func GetClaimsWithAccessToken(log *zap.Logger, token *jwt.Token) (*dto.AccessTokenParse, error) {

	if !token.Valid {
		log.Error("token is not valid")
		return &dto.AccessTokenParse{}, errors.ErrInvalidToken
	}

	claims, ok := token.Claims.(*dto.AccessTokenParse)
	if !ok {
		log.Error("token claims is not found")
		return &dto.AccessTokenParse{}, errors.ErrInvalidToken
	}

	return claims, nil
}

func GetClaimsWithRefreshToken(log *zap.Logger, token *jwt.Token) (*dto.RefreshTokenParse, error) {

	if !token.Valid {
		log.Error("token is not valid")
		return &dto.RefreshTokenParse{}, errors.ErrInvalidToken
	}

	claims, ok := token.Claims.(*dto.RefreshTokenParse)
	if !ok {
		log.Error("token claims is not found")
		return &dto.RefreshTokenParse{}, errors.ErrInvalidToken
	}

	return claims, nil
}
