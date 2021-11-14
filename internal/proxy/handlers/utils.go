package proxy

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gbrlsnchs/jwt/v3"
	"github.com/kenriortega/ngonx/pkg/errors"
	"github.com/kenriortega/ngonx/pkg/logger"
	"github.com/kenriortega/ngonx/pkg/otelify"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func otelRegisterByRequest(ctx context.Context, start time.Time, req *http.Request, err error) {

	traceID := trace.SpanContextFromContext(ctx).TraceID().String()

	otelify.MetricRequestLatencyProxy.(prometheus.ExemplarObserver).ObserveWithExemplar(
		time.Since(start).Seconds(), prometheus.Labels{"traceID": traceID},
	)

	if err != nil {
		logger.LogError(
			"proxy.Director.Metric",
			zap.String("traceID", traceID),
			zap.String("path", req.URL.Path),
			zap.Duration("latency", time.Since(start)),
		)

		return
	}
	logger.LogInfo(
		"proxy.Director.Metric",
		zap.String("traceID", traceID),
		zap.String("path", req.URL.Path),
		zap.Duration("latency", time.Since(start)),
	)
}

// checkJWT check jwt for request
func checkJWT(ctx context.Context, req *http.Request, key string) error {
	ctx, span := otel.Tracer("proxy.gateway.checkJWT").Start(ctx, "checkJWT")
	defer span.End()
	traceID := trace.SpanContextFromContext(ctx).TraceID().String()

	header := req.Header.Get("Authorization") // pass to constanst
	hs := jwt.NewHS256([]byte(key))
	now := time.Now()
	if !strings.HasPrefix(header, "Bearer ") {
		otelify.InstrumentedError(span, "checkJWT.bearer", traceID, errors.ErrBearerTokenFormat)
		return errors.ErrBearerTokenFormat
	}

	token := strings.Split(header, " ")[1]
	pl := JWTPayload{}
	expValidator := jwt.ExpirationTimeValidator(now)
	validatePayload := jwt.ValidatePayload(&pl.Payload, expValidator)

	_, err := jwt.Verify([]byte(token), hs, &pl, validatePayload)

	if errors.ErrorIs(err, jwt.ErrExpValidation) {
		otelify.InstrumentedError(span, "checkJWT.expValidation", traceID, errors.ErrTokenExpValidation)
		return errors.ErrTokenExpValidation
	}
	if errors.ErrorIs(err, jwt.ErrHMACVerification) {
		otelify.InstrumentedError(span, "checkJWT.HMACValidation", traceID, errors.ErrTokenHMACValidation)
		return errors.ErrTokenHMACValidation
	}
	otelify.InstrumentedInfo(span, "checkJWT", traceID)
	return nil
}

// checkAPIKEY check apikey from request
func checkAPIKEY(
	ctx context.Context,
	req *http.Request,
	ph *ProxyHandler,
	engine, key string,
) error {
	ctx, span := otel.Tracer("proxy.gateway.checkAPIKey").Start(ctx, "checkAPIKEY")
	defer span.End()
	traceID := trace.SpanContextFromContext(ctx).TraceID().String()

	header := req.Header.Get("X-API-KEY")
	apikey, err := ph.Service.GetKEY(engine, key)
	if err != nil {
		otelify.InstrumentedError(span, "checkAPIKEY.GetKEY", traceID, errors.ErrGetkeyView)
		return errors.ErrGetkeyView
	}
	if apikey == header {
		otelify.InstrumentedInfo(span, "checkAPIKEY", traceID)
		return nil
	} else {
		invalidKeyErr := errors.NewError("Invalid API KEY")
		otelify.InstrumentedError(span, "chackAPIKEY.invalidHeader", traceID, invalidKeyErr)
		return invalidKeyErr
	}
}
