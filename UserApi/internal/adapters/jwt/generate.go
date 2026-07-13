package jwt

import (
	"ITK_Code/m/v2/internal/core/auth"
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
) (string, error) {

	claimsRefreshToken := auth.RefreshTokenParse{
		AccessTokenJti: jti,
		Jti:            uuid.NewString(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenTTL)),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsRefreshToken)

	refreshTokenString, err := refreshToken.SignedString([]byte(secret))
	if err != nil {
		return "", auth.ErrGenerateToken
	}

	return refreshTokenString, nil
}

func generateAccessToken(
	secret string,
	accessTokenTTL time.Duration,
	user user.User,
	deviceId string,
) (string, string, error) {

	tokenID := uuid.NewString()

	claimsAccessToken := auth.AccessTokenParse{
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
		return "", "", auth.ErrGenerateToken
	}

	return accessTokenString, tokenID, nil
}

func GetClaimsWithAccessToken(log *zap.Logger, token *jwt.Token) (*auth.AccessTokenParse, error) {

	if !token.Valid {
		log.Error("token is not valid")
		return &auth.AccessTokenParse{}, auth.ErrInvalidToken
	}

	claims, ok := token.Claims.(*auth.AccessTokenParse)
	if !ok {
		log.Error("token claims is not found")
		return &auth.AccessTokenParse{}, auth.ErrInvalidToken
	}

	return claims, nil
}

func GetClaimsWithRefreshToken(log *zap.Logger, token *jwt.Token) (*auth.RefreshTokenParse, error) {

	if !token.Valid {
		log.Error("token is not valid")
		return &auth.RefreshTokenParse{}, auth.ErrInvalidToken
	}

	claims, ok := token.Claims.(*auth.RefreshTokenParse)
	if !ok {
		log.Error("token claims is not found")
		return &auth.RefreshTokenParse{}, auth.ErrInvalidToken
	}

	return claims, nil
}
