package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
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
	isCompareHash := strings.Compare(currentValue, hash)
	if isCompareHash != 0 {
		return ErrHashMissMatch
	}
	return nil
}
