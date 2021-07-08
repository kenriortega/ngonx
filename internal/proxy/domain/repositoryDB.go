package proxy

import (
	"fmt"
	"os"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/kenriortega/goproxy/internal/platform/logger"
)

type ProxyRepositoryStorage struct {
	clientBadger *badger.DB
}

func NewProxyRepository(clients ...interface{}) ProxyRepositoryStorage {
	var proxyRepositoryDB ProxyRepositoryStorage
	for _, c := range clients {
		switch c.(type) {
		case *badger.DB:
			proxyRepositoryDB.clientBadger = c.(*badger.DB)
		}
	}
	return proxyRepositoryDB
}

func (r ProxyRepositoryStorage) SaveKEY(engine, key, apikey string) error {
	switch engine {
	case "badger":
		if err := r.clientBadger.Update(func(txn *badger.Txn) error {
			if err := txn.Set([]byte(key), []byte(apikey)); err != nil {
				logger.LogError("savekey: failed")
				return err
			}
			logger.LogInfo("savekey: successful")

			return nil
		}); err != nil {
			logger.LogError("savekey: failed")

			return err
		}

		return nil
	case "local":
		f, err := os.Create("./apikey")

		if err != nil {
			logger.LogError(err.Error())
			return err
		}

		defer f.Close()

		data := []byte(fmt.Sprintf("%s:%s", key, apikey))

		_, err = f.Write(data)

		if err != nil {
			logger.LogError(err.Error())
			return err
		}
		return nil
	}
	return nil
}

func (r ProxyRepositoryStorage) GetKEY(key string) (string, error) {
	var apikey string
	fmt.Println(key)
	if err := r.clientBadger.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		if err := item.Value(func(value []byte) error {
			apikey = string(value)
			return nil
		}); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return "", err
	}

	return apikey, nil
}
