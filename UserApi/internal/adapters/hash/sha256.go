package hash

import (
	"bytes"
	"crypto/sha256"
)

func GenerateHashSHA256(value string) []byte {
	hash := sha256.Sum256([]byte(value))
	return hash[:]
}

func CompareHashSHA256(value string, hash []byte) bool {
	currentValue := GenerateHashSHA256(value)
	isCompareHash := bytes.Compare(currentValue, hash)
	return isCompareHash == 0
}
