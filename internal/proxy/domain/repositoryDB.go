package proxy

import (
	"fmt"
	"os"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/kenriortega/goproxy/pkg/errors"
	"github.com/kenriortega/goproxy/pkg/logger"
)

// ProxyRepositoryStorage struct repository storage
type ProxyRepositoryStorage struct {
	clientBadger *badger.DB
}

// NewProxyRepository return a new ProxyRepositoryStorage
// with a client `*badger.DB`
func NewProxyRepository(clients ...interface{}) ProxyRepositoryStorage {
	var proxyRepositoryDB ProxyRepositoryStorage
	for _, c := range clients {
		switch c := c.(type) {
		case *badger.DB:
			proxyRepositoryDB.clientBadger = c
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
			logger.LogInfo("savekey: successful")

			return nil
		}); err != nil {
			logger.LogError(errors.ErrSavekeyUpdate.Error())
			return errors.ErrSavekeyUpdate
		}

		return nil
	case "local":
		f, err := os.Create("./apikey")

		if err != nil {
			logger.LogError(errors.ErrSavekeyCreateLocal.Error())
			return errors.ErrSavekeyCreateLocal
		}

		defer f.Close()

		data := []byte(fmt.Sprintf("%s:%s", key, apikey))

		_, err = f.Write(data)

		if err != nil {
			logger.LogError(errors.ErrSavekeyWriteOnLocal.Error())
			return errors.ErrSavekeyWriteOnLocal
		}
		return nil
	}
	return nil
}

// GetKEY get key from the database
func (r ProxyRepositoryStorage) GetKEY(key string) (string, error) {
	var apikey string
	fmt.Println(key)
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

	return apikey, nil
}
