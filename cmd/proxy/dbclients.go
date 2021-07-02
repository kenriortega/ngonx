package proxy

import (
	"egosystem.org/micros/internal"
	badger "github.com/dgraph-io/badger/v3"
)

func GetBadgerDB(embedMem bool) *badger.DB {
	var opt badger.Options
	if embedMem {
		opt = badger.DefaultOptions("").WithInMemory(true)
	} else {
		opt = badger.DefaultOptions("./badger.data")
	}

	db, err := badger.Open(opt)
	if err != nil {
		internal.LogError(err.Error())
	}
	// defer db.Close()
	return db
}
