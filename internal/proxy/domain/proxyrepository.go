package proxy

import (
	"context"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/go-redis/redis/v8"
	"github.com/kenriortega/ngonx/pkg/errors"
	"github.com/kenriortega/ngonx/pkg/otelify"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// ProxyRepositoryStorage struct repository storage
type ProxyRepositoryStorage struct {
	clientBadger *badger.DB
	clientRdb    *redis.Client
}

// NewProxyRepository return a new ProxyRepositoryStorage
// with a client `*badger.DB`
func NewProxyRepository(clients ...interface{}) ProxyRepositoryStorage {
	var proxyRepositoryDB ProxyRepositoryStorage
	for _, c := range clients {
		switch c := c.(type) {
		case *badger.DB:
			proxyRepositoryDB.clientBadger = c
		case *redis.Client:
			proxyRepositoryDB.clientRdb = c
		}
	}
	return proxyRepositoryDB
}

// SaveKEY save a key on the database
func (r ProxyRepositoryStorage) SaveKEY(engine, key, apikey string) error {
	ctx, span := otel.Tracer("proxy.repo").Start(context.Background(), "SaveKEY")
	defer span.End()
	traceID := trace.SpanContextFromContext(ctx).TraceID().String()
	switch engine {
	case "badger":
		if err := r.clientBadger.Update(func(txn *badger.Txn) error {
			if err := txn.Set([]byte(key), []byte(apikey)); err != nil {
				otelify.InstrumentedError(span, "badger", traceID, err)
				return errors.ErrSavekeyUpdateTX
			}
			otelify.InstrumentedInfo(span, "repo.SaveKey", traceID)
			return nil
		}); err != nil {
			otelify.InstrumentedError(span, "badger", traceID, err)
			return errors.ErrSavekeyUpdate
		}

		return nil
	case "redis":
		if _, err := r.clientRdb.HSet(context.TODO(), key, apikey).Result(); err != nil {
			otelify.InstrumentedError(span, "redis", traceID, err)
		}
		otelify.InstrumentedInfo(span, "repo.SaveKey", traceID)
	}
	return nil
}

// GetKEY get key from the database
func (r ProxyRepositoryStorage) GetKEY(engine, key string) (string, error) {
	ctx, span := otel.Tracer("proxy.repo").Start(context.Background(), "GetKEY")
	defer span.End()
	traceID := trace.SpanContextFromContext(ctx).TraceID().String()
	var apikey string

	switch engine {
	case "badger":
		if err := r.clientBadger.View(func(txn *badger.Txn) error {
			item, err := txn.Get([]byte(key))
			if err != nil {
				otelify.InstrumentedError(span, "badger", traceID, err)
				return errors.ErrGetkeyTX
			}
			if err := item.Value(func(value []byte) error {
				apikey = string(value)
				return nil
			}); err != nil {
				otelify.InstrumentedError(span, "badger", traceID, err)
				return errors.ErrGetkeyValue
			}

			return nil
		}); err != nil {
			otelify.InstrumentedError(span, "badger", traceID, err)
			return "", errors.ErrGetkeyView
		}
	case "redis":
		value, err := r.clientRdb.Get(context.TODO(), key).Result()
		if err == redis.Nil || err != nil {
			otelify.InstrumentedError(span, "redis", traceID, err)
			return "", err
		}
		apikey = value
	}
	otelify.InstrumentedInfo(span, "repo.GetKey", traceID)

	return apikey, nil
}
