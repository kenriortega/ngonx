package internal

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"time"
)

var (
	seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
)

const CHARSET = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789$*-/+"

func ApiKeyGenerator(word string) string {
	h := sha256.New()
	_, err := h.Write([]byte(word))
	if err != nil {
		LogError("Error")
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

// StringWithCharset generate words
func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
