package gateway

import (
	"fmt"
	"os"

	"egosystem.org/micros/internal"
	badger "github.com/dgraph-io/badger/v3"
)

type ProxyRepositoryDB struct {
	clientBadger *badger.DB
}

func NewProxyRepository(clients ...interface{}) ProxyRepositoryDB {
	var proxyRepositoryDB ProxyRepositoryDB
	for _, c := range clients {
		switch c.(type) {
		case *badger.DB:
			proxyRepositoryDB.clientBadger = c.(*badger.DB)
		}
	}
	return proxyRepositoryDB
}

func (r ProxyRepositoryDB) SaveKEY(engine, apikey string) error {
	switch engine {
	case "badger":
		err := r.clientBadger.Update(func(txn *badger.Txn) error {
			err := txn.Set([]byte("apikey"), []byte(apikey))
			internal.LogError(err.Error())
			return err
		})
		if err != nil {
			internal.LogError(err.Error())
			return err
		}
		defer r.clientBadger.Close()
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
