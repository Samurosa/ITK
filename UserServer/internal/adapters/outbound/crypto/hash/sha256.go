package hash

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
)

var (
	ErrHashMissMatch = errors.New("the hash does not match")
)

type SHA256 struct{}

func (SHA256) GenerateHashSHA256(value string) string { return generateHash(value) }

func (SHA256) CompareHashSHA256(value string, hash string) error {
	currentValue := generateHash(value)

	if subtle.ConstantTimeCompare([]byte(currentValue), []byte(hash)) != 1 {
		return ErrHashMissMatch
	}

	return nil
}

func generateHash(value string) string {
	hash := sha256.Sum256([]byte(value))

	hashToString := hex.EncodeToString(hash[:])
	return hashToString
}
