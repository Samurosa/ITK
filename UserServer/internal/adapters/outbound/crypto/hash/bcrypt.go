package hash

import (
	"ITK_Code/m/v2/internal/core/auth"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Bcrypt struct{}

func (Bcrypt) GeneratePasswordHash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), 14)
}

func (Bcrypt) VerifyPasswordHash(password string, hash []byte) error {
	err := bcrypt.CompareHashAndPassword(hash, []byte(password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return auth.ErrIncorrectPassword
	}
	return err
}
