package badgerdb

import (
	badger "github.com/dgraph-io/badger/v3"
	"github.com/kenriortega/ngonx/pkg/logger"
)

var pathDB = "./badger.data"

// GetBadgerDB return `*badger.DB`
// this client provide GET and SAVE methods
func GetBadgerDB(embedMem bool) *badger.DB {
	var opt badger.Options
	var clientBadger *badger.DB
	if embedMem {
		opt = badger.DefaultOptions("").WithInMemory(true).WithBypassLockGuard(true)
	} else {
		opt = badger.DefaultOptions(pathDB).WithBypassLockGuard(true)
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
