package jwt

import (
	"ITK_Code/m/v2/internal/core/corerrors"
	"ITK_Code/m/v2/internal/core/user"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func generateRefreshToken(
	secret string,
	refreshTokenTTL time.Duration,
	jti string,
) (string, RefreshToken, error) {

	claimsRefreshToken := RefreshToken{
		AccessTokenJTI:  jti,
		RefreshTokenJTI: uuid.NewString(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenTTL)),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsRefreshToken)

	refreshTokenString, err := refreshToken.SignedString([]byte(secret))
	if err != nil {
		return "", RefreshToken{}, corerrors.ErrGenerateToken
	}

	return refreshTokenString, claimsRefreshToken, nil
}

func generateAccessToken(
	secret string,
	accessTokenTTL time.Duration,
	user user.User,
	deviceId string,
) (string, AccessToken, error) {

	tokenID := uuid.NewString()
	issuedAt := time.Now()

	claimsAccessToken := AccessToken{
		UserID:   user.ID,
		Role:     string(user.Role),
		DeviceID: deviceId,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID,
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenTTL)),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsAccessToken)

	accessTokenString, err := accessToken.SignedString([]byte(secret))
	if err != nil {
		return "", AccessToken{}, corerrors.ErrGenerateToken
	}

	return accessTokenString, claimsAccessToken, nil
}

func GetClaimsWithRefreshToken(token *jwt.Token) (*RefreshToken, error) {

	if !token.Valid {
		return &RefreshToken{}, corerrors.ErrInvalidToken
	}

	claims, ok := token.Claims.(*RefreshToken)
	if !ok {
		return &RefreshToken{}, corerrors.ErrInvalidToken
	}

	return &RefreshToken{
		AccessTokenJTI:  claims.AccessTokenJTI,
		RefreshTokenJTI: claims.RefreshTokenJTI,
	}, nil
}
