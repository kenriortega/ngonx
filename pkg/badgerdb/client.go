package badgerdb

import (
	"context"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/kenriortega/ngonx/pkg/errors"
	"github.com/kenriortega/ngonx/pkg/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
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
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		logger.LogError(
			errors.Errorf("badger: %v", err).Error(),
			zap.String("traceID", traceID),
		)

		panic(err)
	}
	clientBadger = db
	logger.LogInfo(
		"proxy.GetBadgerDB",
		zap.String("traceID", traceID),
	)
	// defer clientBadger.Close()
	span.SetAttributes(attribute.String("badgerdb.create.client", "Success"))
	return clientBadger
}
