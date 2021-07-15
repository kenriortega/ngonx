package proxy

import (
	"github.com/kenriortega/goproxy/internal/platform/badgerdb"
	domain "github.com/kenriortega/goproxy/internal/proxy/domain"

	"testing"
)

var proxyRepository domain.ProxyRepository

func Test_SaveSecretKEY(t *testing.T) {
	clientBadger := badgerdb.GetBadgerDB(true)
	proxyRepository = domain.NewProxyRepository(clientBadger)
	err := proxyRepository.SaveKEY("badger", "key", "apikey")
	if err != nil {
		t.Error("Error to created key")
	}
}

func Test_GetKEY(t *testing.T) {
	clientBadger := badgerdb.GetBadgerDB(true)
	proxyRepository = domain.NewProxyRepository(clientBadger)
	result, err := proxyRepository.GetKEY("key")
	if err != nil {
		t.Error("Error to created key")
	}
	if result != "apikey" {
		t.Error("Expected apikey as a value for `key`")
	}
}
