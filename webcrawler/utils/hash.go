package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func GenerateContentHash(content string) string {
	hasher := sha256.New()
	hasher.Write([]byte(content))
	return hex.EncodeToString(hasher.Sum(nil))
}
