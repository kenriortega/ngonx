package badgerdb

import (
	"context"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/kenriortega/ngonx/pkg/otelify"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var pathDB = "./badger.data"

// GetBadgerDB return `*badger.DB`
// this client provide GET and SAVE methods
func GetBadgerDB(ctx context.Context, embedMem bool) *badger.DB {
	ctx, span := otel.Tracer("badger.client").Start(ctx, "GetBadgerDB")
	defer span.End()
	traceID := trace.SpanContextFromContext(ctx).TraceID().String()

	var opt badger.Options
	var clientBadger *badger.DB
	if embedMem {
		opt = badger.DefaultOptions("").WithInMemory(true).WithBypassLockGuard(true)
	} else {
		opt = badger.DefaultOptions(pathDB).WithBypassLockGuard(true)
	}

	db, err := badger.Open(opt)
	if err != nil {
		otelify.InstrumentedError(span, "badger", traceID, err)
		panic(err)
	}
	clientBadger = db

	otelify.InstrumentedInfo(span, "proxy.GetBadgerDB", traceID)

	return clientBadger
}
