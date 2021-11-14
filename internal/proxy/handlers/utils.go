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
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func otelRegister(ctx context.Context, start time.Time, req *http.Request, err error) {

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
	_, span := otel.Tracer("proxy.gateway.checkJWT").Start(ctx, "checkJWT")
	defer span.End()
	header := req.Header.Get("Authorization") // pass to constanst
	hs := jwt.NewHS256([]byte(key))
	now := time.Now()
	if !strings.HasPrefix(header, "Bearer ") {
		span.RecordError(errors.ErrBearerTokenFormat)
		span.SetStatus(codes.Error, errors.ErrBearerTokenFormat.Error())
		logger.LogError(errors.Errorf("proxy: %v", errors.ErrBearerTokenFormat).Error())

		return errors.ErrBearerTokenFormat
	}

	token := strings.Split(header, " ")[1]
	pl := JWTPayload{}
	expValidator := jwt.ExpirationTimeValidator(now)
	validatePayload := jwt.ValidatePayload(&pl.Payload, expValidator)

	_, err := jwt.Verify([]byte(token), hs, &pl, validatePayload)

	if errors.ErrorIs(err, jwt.ErrExpValidation) {
		span.RecordError(errors.ErrTokenExpValidation)
		span.SetStatus(codes.Error, errors.ErrTokenExpValidation.Error())
		logger.LogError(errors.Errorf("proxy: %v", errors.ErrTokenExpValidation).Error())

		return errors.ErrTokenExpValidation
	}
	if errors.ErrorIs(err, jwt.ErrHMACVerification) {
		span.RecordError(errors.ErrTokenHMACValidation)
		span.SetStatus(codes.Error, errors.ErrTokenHMACValidation.Error())
		logger.LogError(errors.Errorf("proxy: %v", errors.ErrTokenHMACValidation).Error())

		return errors.ErrTokenHMACValidation
	}
	span.AddEvent("checkJWT done!")
	return nil
}

// checkAPIKEY check apikey from request
func checkAPIKEY(
	ctx context.Context,
	req *http.Request,
	ph *ProxyHandler,
	engine, key string,
) error {
	_, span := otel.Tracer("proxy.gateway.checkAPIKey").Start(ctx, "checkAPIKEY")
	defer span.End()
	apikey, err := ph.Service.GetKEY(engine, key)
	header := req.Header.Get("X-API-KEY")
	if err != nil {
		span.RecordError(errors.ErrGetkeyView)
		span.SetStatus(codes.Error, errors.ErrGetkeyView.Error())
		logger.LogError(errors.Errorf("proxy: %v", errors.ErrGetkeyView).Error())

	}
	if apikey == header {
		logger.LogInfo("proxy: check secret from request OK")
		return nil
	} else {
		invalidKeyErr := errors.NewError("Invalid API KEY")
		span.RecordError(invalidKeyErr)
		span.SetStatus(codes.Error, invalidKeyErr.Error())
		logger.LogError(invalidKeyErr.Error())
		return invalidKeyErr
	}

}
