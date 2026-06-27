package jwt

import (
	"ITK_Code/m/v2/internal/core/auth"
	"ITK_Code/m/v2/internal/core/user/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func generateRefreshToken() (string, error) {
	refreshTokenString := uuid.NewString()
	if refreshTokenString == "" {
		return "", auth.ErrGenerateRefreshToken
	}

	return refreshTokenString, nil
}

func generateAccessToken(
	secret string,
	accessTokenTTL time.Duration,
	user models.User,
) (string, error) {
	accessToken := jwt.New(jwt.SigningMethodHS256)

	claimsAccessToken := accessToken.Claims.(jwt.MapClaims)
	claimsAccessToken["id"] = user.ID
	claimsAccessToken["role"] = user.Role
	claimsAccessToken["jti"] = uuid.NewString()
	claimsAccessToken["exp"] = time.Now().Add(accessTokenTTL).Unix()

	accessTokenString, err := accessToken.SignedString(secret)
	if err != nil {
		return "", auth.ErrGenerateAccessToken
	}

	return accessTokenString, nil
}
