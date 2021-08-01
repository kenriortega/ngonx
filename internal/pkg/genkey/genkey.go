package genkey

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"time"

	"github.com/kenriortega/goproxy/internal/pkg/errors"
	"github.com/kenriortega/goproxy/internal/pkg/logger"
)

var (
	// seededRand random number
	seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
)

const (
	// CHARSET of characters
	CHARSET = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789$*-/+"
	// SALT for define the size of the buffer
	SALT = 10
)

// ApiKeyGenerator it`s used to create a hash from a random string
//
// Example
// input word 1q2w3e4r5t the expected result are 28f0116ef42bf718324946f13d787a1d41274a08335d52ee833d5b577f02a32a
func ApiKeyGenerator(word string) string {
	h := sha256.New()
	_, err := h.Write([]byte(word))
	if err != nil {
		logger.LogError(errors.ErrApiKeyGenerator.Error())
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

// StringWithCharset generate randoms words
func StringWithCharset() string {
	b := make([]byte, SALT)
	for i := range b {
		b[i] = CHARSET[seededRand.Intn(len(CHARSET))]
	}
	return string(b)
}
