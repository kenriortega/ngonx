package proxy

import (
	"context"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/go-redis/redis/v8"
	"github.com/kenriortega/ngonx/pkg/errors"
	"github.com/kenriortega/ngonx/pkg/logger"
)

// ProxyRepositoryStorage struct repository storage
type ProxyRepositoryStorage struct {
	clientBadger *badger.DB
	clientRdb    *redis.Client
}

// NewProxyRepository return a new ProxyRepositoryStorage
// with a client `*badger.DB`
func NewProxyRepository(clients ...interface{}) ProxyRepositoryStorage {
	var proxyRepositoryDB ProxyRepositoryStorage
	for _, c := range clients {
		switch c := c.(type) {
		case *badger.DB:
			proxyRepositoryDB.clientBadger = c
		case *redis.Client:
			proxyRepositoryDB.clientRdb = c
		}
	}
	return proxyRepositoryDB
}

// SaveKEY save a key on the database
func (r ProxyRepositoryStorage) SaveKEY(engine, key, apikey string) error {
	switch engine {
	case "badger":
		if err := r.clientBadger.Update(func(txn *badger.Txn) error {
			if err := txn.Set([]byte(key), []byte(apikey)); err != nil {
				logger.LogError(errors.ErrSavekeyUpdateTX.Error())
				return errors.ErrSavekeyUpdateTX
			}
			logger.LogInfo("proxy: savekey was successful")

			return nil
		}); err != nil {

			return errors.ErrSavekeyUpdate
		}

		return nil
	case "redis":
		if _, err := r.clientRdb.HSet(context.TODO(), key, apikey).Result(); err != nil {
			logger.LogError(errors.Errorf("proxy redis: %v", err).Error())
		}
	}
	return nil
}

// GetKEY get key from the database
func (r ProxyRepositoryStorage) GetKEY(engine, key string) (string, error) {
	var apikey string

	switch engine {
	case "badger":
		if err := r.clientBadger.View(func(txn *badger.Txn) error {
			item, err := txn.Get([]byte(key))
			if err != nil {
				return errors.ErrGetkeyTX
			}
			if err := item.Value(func(value []byte) error {
				apikey = string(value)
				return nil
			}); err != nil {
				return errors.ErrGetkeyValue
			}

			return nil
		}); err != nil {
			return "", errors.ErrGetkeyView
		}
	case "redis":
		value, err := r.clientRdb.Get(context.TODO(), key).Result()
		if err == redis.Nil || err != nil {
			return "", err
		}
		apikey = value
	}

	return apikey, nil
}
