package hash

import "crypto/sha256"

func GenerateHashSHA256(value string) [32]byte {
	return sha256.Sum256([]byte(value))
}

func CompareHashSHA256(value string, hash [32]byte) bool {
	currentValue := GenerateHashSHA256(value)
	return currentValue == hash
}
