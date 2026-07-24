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

func GenerateHashSHA256(value string) string {
	hash := sha256.Sum256([]byte(value))

	hashToString := hex.EncodeToString(hash[:])
	return hashToString
}

func CompareHashSHA256(value string, hash string) error {
	currentValue := GenerateHashSHA256(value)

	if subtle.ConstantTimeCompare([]byte(currentValue), []byte(hash)) != 1 {
		return ErrHashMissMatch
	}

	return nil
}
