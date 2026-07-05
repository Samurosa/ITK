package jwt

import (
	"ITK_Code/m/v2/internal/core/auth"
	"ITK_Code/m/v2/internal/core/user"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func generateRefreshToken() (string, error) {
	refreshTokenString := uuid.NewString()

	return refreshTokenString, nil
}

func generateAccessToken(
	secret string,
	accessTokenTTL time.Duration,
	user user.User,
) (string, error) {

	claimsAccessToken := auth.TokenParse{
		UserID: user.ID,
		Role:   string(user.Role),
		Jti:    uuid.NewString(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenTTL)),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsAccessToken)

	accessTokenString, err := accessToken.SignedString([]byte(secret))
	if err != nil {
		return "", auth.ErrGenerateToken
	}

	return accessTokenString, nil
}

func GetClaimsWithToken(log *zap.Logger, token *jwt.Token) (*auth.TokenParse, error) {

	if !token.Valid {
		log.Error("token is not valid")
		return &auth.TokenParse{}, auth.ErrInvalidToken
	}

	claims, ok := token.Claims.(*auth.TokenParse)
	if !ok {
		log.Error("token claims is not found")
		return &auth.TokenParse{}, auth.ErrInvalidToken
	}

	return claims, nil
}
