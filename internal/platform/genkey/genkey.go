package genkey

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"time"

	"github.com/kenriortega/goproxy/internal/platform/errors"
	"github.com/kenriortega/goproxy/internal/platform/logger"
)

var (
	seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
)

const (
	CHARSET = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789$*-/+"
	SALT = 10
)

func ApiKeyGenerator(word string) string {
	h := sha256.New()
	_, err := h.Write([]byte(word))
	if err != nil {
		logger.LogError(errors.ErrApiKeyGenerator.Error())
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

// StringWithCharset generate words
func StringWithCharset() string {
	b := make([]byte, SALT)
	for i := range b {
		b[i] = CHARSET[seededRand.Intn(len(CHARSET))]
	}
	return string(b)
}
