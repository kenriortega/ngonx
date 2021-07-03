package gateway

import (
	"fmt"
	"os"

	"egosystem.org/micros/internal"
	badger "github.com/dgraph-io/badger/v3"
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
				internal.LogError("savekey: failed")
				return err
			}
			internal.LogInfo("savekey: successful")

			return nil
		}); err != nil {
			internal.LogError("savekey: failed")

			return err
		}

		return nil
	case "local":
		f, err := os.Create("./apikey")

		if err != nil {
			internal.LogError(err.Error())
			return err
		}

		defer f.Close()

		data := []byte(fmt.Sprintf("apikey:%s", apikey))

		_, err = f.Write(data)

		if err != nil {
			internal.LogError(err.Error())
			return err
		}
		return nil
	}
	return nil
}

func (r ProxyRepositoryStorage) GetKEY(key string) (string, error) {
	var apikey string

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
