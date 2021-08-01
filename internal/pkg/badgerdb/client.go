package badgerdb

import (
	badger "github.com/dgraph-io/badger/v3"
	"github.com/kenriortega/goproxy/internal/pkg/logger"
)

func GetBadgerDB(embedMem bool) *badger.DB {
	var opt badger.Options
	var clientBadger *badger.DB
	if embedMem {
		opt = badger.DefaultOptions("").WithInMemory(true).WithBypassLockGuard(true)
	} else {
		opt = badger.DefaultOptions("./badger.data").WithBypassLockGuard(true)
	}

	db, err := badger.Open(opt)
	if err != nil {
		logger.LogError(err.Error())

		panic(err)
	}
	clientBadger = db
	// defer clientBadger.Close()

	return clientBadger
}
