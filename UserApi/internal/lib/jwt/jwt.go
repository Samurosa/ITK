package jwt

import (
	"ITK_Code/m/v2/internal/domain/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func NewToken(user models.User, secret models.App, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims[user.ID] = user.ID
	claims[user.Login] = user.Login
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
