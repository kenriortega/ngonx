package proxy

import (
	"context"

	domain "github.com/kenriortega/ngonx/internal/proxy/domain"
	"github.com/kenriortega/ngonx/pkg/otelify"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// ProxyService interface service for proxy repository funcionalities
type ProxyService interface {
	SaveSecretKEY(string, string, string) error
	GetKEY(string, string) (string, error)
}

// DefaultProxyService struct for management proxy repository
type DefaultProxyService struct {
	repo domain.ProxyRepository
}

// NewProxyService return new DefaultProxyService
func NewProxyService(repository domain.ProxyRepository) DefaultProxyService {
	return DefaultProxyService{repo: repository}
}

// SaveSecretKEY save secret key
func (s DefaultProxyService) SaveSecretKEY(engine, key, apikey string) (string, error) {
	ctx, span := otel.Tracer("proxy.service.SaveSecretKEY").Start(context.Background(), "ProxyGateway")
	defer span.End()
	traceID := trace.SpanContextFromContext(ctx).TraceID().String()
	err := s.repo.SaveKEY(engine, key, apikey)
	if err != nil {
		otelify.InstrumentedError(span, "SaveKey", traceID, err)
		return "failed", err
	}
	otelify.InstrumentedInfo(span, "service.SaveSecretKEY", traceID)
	return "ok", nil
}

// GetKEY get key
func (s DefaultProxyService) GetKEY(engine, key string) (string, error) {
	ctx, span := otel.Tracer("proxy.service.GetKEY").Start(context.Background(), "ProxyGateway")
	defer span.End()
	traceID := trace.SpanContextFromContext(ctx).TraceID().String()
	result, err := s.repo.GetKEY(engine, key)
	if err != nil {
		otelify.InstrumentedError(span, "GetKey", traceID, err)
		return "failed", err
	}
	otelify.InstrumentedInfo(span, "service.GetKey", traceID)
	return result, nil
}
