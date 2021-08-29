package mngt

import (
	"fmt"

	"github.com/dgraph-io/badger/v3"
	"github.com/go-redis/redis/v8"
)

var endpoints []Endpoint

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

	return endpoints, nil
}

func (r MngtRepositoryStorage) RegisterEnpoint(endpoint Endpoint) error {

	endpoints = append(endpoints, endpoint)
	fmt.Println(endpoints)
	return nil
}
