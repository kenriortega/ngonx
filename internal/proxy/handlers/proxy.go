package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/kenriortega/goproxy/pkg/errors"
	"github.com/kenriortega/goproxy/pkg/logger"

	"github.com/gbrlsnchs/jwt/v3"
	domain "github.com/kenriortega/goproxy/internal/proxy/domain"
	services "github.com/kenriortega/goproxy/internal/proxy/services"
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
		logger.LogInfo(result)
	}
	logger.LogInfo(result)
}

// ProxyGateway handler for management all request
func (ph *ProxyHandler) ProxyGateway(endpoints domain.ProxyEndpoint, key, securityType string) {
	for _, endpoint := range endpoints.Endpoints {

		target, err := url.Parse(
			fmt.Sprintf("%s%s", endpoints.HostURI, endpoint.PathEndpoint),
		)
		if err != nil {
			logger.LogError(err.Error())
		}
		if endpoint.PathProtected {
			proxy = httputil.NewSingleHostReverseProxy(target)

			originalDirector := proxy.Director
			proxy.Director = func(req *http.Request) {
				originalDirector(req)

				switch securityType {
				case "jwt":
					err := checkJWTSecretKeyFromRequest(req, key)
					proxy.ModifyResponse = modifyResponse(err)
				case "apikey":

					err := checkAPIKEYSecretKeyFromRequest(req, ph, key)
					proxy.ModifyResponse = modifyResponse(err)
				}

			}
			proxy.ErrorHandler = func(rw http.ResponseWriter, r *http.Request, err error) {
				rw.WriteHeader(http.StatusInternalServerError)
				_, _ = rw.Write([]byte(err.Error()))
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
}

// checkJWTSecretKeyFromRequest check jwt for request
func checkJWTSecretKeyFromRequest(req *http.Request, key string) error {
	header := req.Header.Get("Authorization") // pass to constanst
	hs := jwt.NewHS256([]byte(key))
	now := time.Now()
	if !strings.HasPrefix(header, "Bearer ") {
		logger.LogError(errors.ErrBearerTokenFormat.Error())
		return errors.ErrBearerTokenFormat
	}
	token := strings.Split(header, " ")[1]
	pl := JWTPayload{}
	expValidator := jwt.ExpirationTimeValidator(now)
	validatePayload := jwt.ValidatePayload(&pl.Payload, expValidator)
	_, err := jwt.Verify([]byte(token), hs, &pl, validatePayload)

	if errors.ErrorIs(err, jwt.ErrExpValidation) {
		logger.LogError(errors.ErrTokenExpValidation.Error())
		return errors.ErrTokenExpValidation
	}
	if errors.ErrorIs(err, jwt.ErrHMACVerification) {
		logger.LogError(errors.ErrTokenHMACValidation.Error())
		return errors.ErrTokenHMACValidation
	}

	return nil
}

// checkAPIKEYSecretKeyFromRequest check apikey from request
func checkAPIKEYSecretKeyFromRequest(req *http.Request, ph *ProxyHandler, key string) error {
	apikey, err := ph.Service.GetKEY(key)
	header := req.Header.Get("X-API-KEY") // pass to constants
	if err != nil {
		logger.LogError(errors.ErrGetkeyView.Error())
	}
	if apikey == header {
		logger.LogInfo("OK")
		return nil
	} else {
		logger.LogInfo("Invalid apikey")
		return errors.NewError("Invalid API KEY")
	}
}

// modifyResponse modify response
func modifyResponse(err error) func(*http.Response) error {
	return func(resp *http.Response) error {
		resp.Header.Set("X-Proxy", "EgoProxy")
		if err != nil {
			return err
		}
		return nil
	}
}
