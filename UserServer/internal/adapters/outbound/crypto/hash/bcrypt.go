package hash

import "golang.org/x/crypto/bcrypt"

type Bcrypt struct{}

func (Bcrypt) GeneratePasswordHash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), 14)
}

func (Bcrypt) VerifyPasswordHash(password string, hash []byte) error {
	return bcrypt.CompareHashAndPassword(hash, []byte(password))
}
