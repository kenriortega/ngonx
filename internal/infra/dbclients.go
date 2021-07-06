package infra

import (
	badger "github.com/dgraph-io/badger/v3"
	"github.com/kenriortega/goproxy/internal/utils"
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
		utils.LogError(err.Error())
	}
	// defer db.Close()
	return db
}
