package utils

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"os"
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

func CreateSettingFile(setingFile string) {
	f, err := os.Create(fmt.Sprintf("./%s", setingFile))
	ymldata := `
proxy:
  host_proxy: 0.0.0.0
  port_proxy: 5000
  cache_proxy:
    engine: badger # local|badgerDB|redis
    key: secretKey
  security:
    type: jwt # apikey|jwt|none
    secret_key: key00 # apikey jwtkey this value can be replace by genkey command
  # maps of microservices with routes
  services_proxy:
      - name: microA
        host_uri: http://localhost:3000
        endpoints:
          - path_endpoints: /api/v1/health/
            path_proxy: /health/
            path_protected: false
		`
	if err != nil {
		logger.LogError(errors.ErrCreatingSettingFile.Error())
	}

	defer f.Close()

	data := []byte(ymldata)

	_, err = f.Write(data)

	if err != nil {
		logger.LogError(errors.ErrWritingSettingFile.Error())
	}
}