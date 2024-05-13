package util

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashValue(value []byte) string {
	h := sha256.New()

	h.Write(value)

	return hex.EncodeToString(h.Sum(nil))
}
