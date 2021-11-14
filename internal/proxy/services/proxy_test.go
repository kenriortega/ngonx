package proxy

import (
	"context"

	domain "github.com/kenriortega/ngonx/internal/proxy/domain"
	"github.com/kenriortega/ngonx/pkg/badgerdb"

	"testing"
)

var proxyRepository domain.ProxyRepository

func Test_SaveSecretKEY(t *testing.T) {
	clientBadger := badgerdb.GetBadgerDB(context.Background(), false)
	proxyRepository = domain.NewProxyRepository(clientBadger)
	err := proxyRepository.SaveKEY("badger", "key", "apikey")
	if err != nil {
		t.Error("Error to created key")
	}
}

func Test_GetKEY(t *testing.T) {
	clientBadger := badgerdb.GetBadgerDB(context.Background(), false)
	proxyRepository = domain.NewProxyRepository(clientBadger)
	result, err := proxyRepository.GetKEY("badger", "key")
	if err != nil {
		t.Error("Error to created key")
	}
	if result != "apikey" {
		t.Error("Expected apikey as a value for `key`")
	}
}
