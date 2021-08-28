package mngt

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/go-redis/redis/v8"
)

type MngtRepositoryStorage struct {
	clientBadger *badger.DB
	clientRdb    *redis.Client
}

func NewMngtRepositoryStorage(clients ...interface{}) MngtRepositoryStorage {
	var mngtRepository MngtRepositoryStorage
	for _, c := range clients {
		switch c := c.(type) {
		case *badger.DB:
			mngtRepository.clientBadger = c
		case *redis.Client:
			mngtRepository.clientRdb = c
		}
	}
	return mngtRepository
}

func (r MngtRepositoryStorage) ListEnpoints() ([]Endpoint, error) {
	var endpoints []Endpoint
	return endpoints, nil
}
