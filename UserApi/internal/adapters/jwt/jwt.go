package jwt

import (
	"ITK_Code/m/v2/internal/core/user/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func NewToken(user models.User, secret models.App, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims[user.Login] = user.Login
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString(secret.Secret)
	if err != nil {
		return "token.SignedString, invalid format token", err
	}

	return tokenString, nil
}

//исправить
