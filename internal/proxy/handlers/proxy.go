package proxy

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/kenriortega/ngonx/pkg/errors"
	"github.com/kenriortega/ngonx/pkg/logger"
	"github.com/kenriortega/ngonx/pkg/otelify"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/gbrlsnchs/jwt/v3"
	domain "github.com/kenriortega/ngonx/internal/proxy/domain"
	services "github.com/kenriortega/ngonx/internal/proxy/services"
)

// proxy global var for management of reverse proxy
var proxy *httputil.ReverseProxy

// JWTPayload custom struc for jwt Payload
type JWTPayload struct {
	jwt.Payload
}

// ResponseMiddleware struct for middleware responses
type ResponseMiddleware struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// ProxyHandler handler for proxy funcionalities
type ProxyHandler struct {
	Service services.DefaultProxyService
}

// SaveSecretKEY handler for save secrets
func (ph *ProxyHandler) SaveSecretKEY(engine, key, apikey string) {
	result, err := ph.Service.SaveSecretKEY(engine, key, apikey)
	if err != nil {
		logger.LogError(errors.Errorf("proxy: SaveSecretKEY %v", err).Error())
	}
	logger.LogInfo("proxy: SaveSecretKEY" + result)
}

// ProxyGateway handler for management all request
func (ph *ProxyHandler) ProxyGateway(
	endpoints domain.ProxyEndpoint,
	engine,
	key,
	securityType string,
) {
	ctx, span := otel.Tracer("proxy.gateway").Start(context.Background(), "ProxyGateway")
	defer span.End()
	traceID := trace.SpanContextFromContext(ctx).TraceID().String()
	for _, endpoint := range endpoints.Endpoints {
		start := time.Now()

		target, err := url.Parse(
			fmt.Sprintf("%s%s", endpoints.HostURI, endpoint.PathEndpoint),
		)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			logger.LogError(errors.Errorf("proxy: %v", err).Error())
		}

		if endpoint.PathProtected {
			var err error
			proxy = httputil.NewSingleHostReverseProxy(target)

			originalDirector := proxy.Director
			proxy.Director = func(req *http.Request) {
				originalDirector(req)
				switch securityType {
				case "jwt":
					err = checkJWT(ctx, req, key)
				case "apikey":
					err = checkAPIKEY(ctx, req, ph, engine, key)
				}
				otelRegisterByRequest(ctx, start, req, err)
			}
			proxy.ModifyResponse = func(resp *http.Response) error {
				resp.Header.Set("X-Proxy", "Ngonx")
				if err != nil {
					return err
				}
				return nil
			}
			proxy.ErrorHandler = func(w http.ResponseWriter, req *http.Request, err error) {
				rpm := ResponseMiddleware{
					Message: err.Error(),
					Code:    http.StatusBadGateway,
				}
				w.WriteHeader(rpm.Code)
				w.Header().Set("Content-Type", "application/json")
				bytes, err := json.Marshal(&rpm)
				if err != nil {
					logger.LogError(err.Error())
				}
				_, err = w.Write(bytes)
				if err != nil {
					logger.LogError(err.Error())
				}

			}
			http.Handle(
				endpoint.PathToProxy,
				http.StripPrefix(
					endpoint.PathToProxy,
					proxy,
				),
			)
		} else {

			proxy = httputil.NewSingleHostReverseProxy(target)

			originalDirector := proxy.Director
			proxy.Director = func(req *http.Request) {
				originalDirector(req)
				otelRegisterByRequest(ctx, start, req, nil)
			}
			proxy.ModifyResponse = func(resp *http.Response) error {
				resp.Header.Set("X-Proxy", "Ngonx")
				return nil
			}

			http.Handle(
				endpoint.PathToProxy,
				http.StripPrefix(
					endpoint.PathToProxy,
					proxy,
				),
			)
		}
	}
	otelify.InstrumentedInfo(span, "proxy.Gateway", traceID)
}
