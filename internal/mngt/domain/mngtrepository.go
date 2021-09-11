package mngt

import (
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

func (r MngtRepositoryStorage) ListEndpoints() ([]Endpoint, error) {

	return endpoints, nil
}

func (r MngtRepositoryStorage) RegisterEndpoint(endpoint Endpoint) error {

	endpoints = append(endpoints, endpoint)

	return nil
}

func (r MngtRepositoryStorage) UpdateEndpoint(endpoint Endpoint) error {

	for idx, it := range endpoints {
		if it.ID == endpoint.ID {
			endpoints[idx].Status = endpoint.Status
		}
	}

	return nil
}
