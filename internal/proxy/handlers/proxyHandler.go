package proxy

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/kenriortega/goproxy/internal/platform/errors"
	"github.com/kenriortega/goproxy/internal/platform/logger"

	"github.com/gbrlsnchs/jwt/v3"
	domain "github.com/kenriortega/goproxy/internal/proxy/domain"
	services "github.com/kenriortega/goproxy/internal/proxy/services"
)

var proxy *httputil.ReverseProxy

type JWTPayload struct {
	jwt.Payload
}
type ResponseMiddleware struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}
type ProxyHandler struct {
	Service services.DefaultProxyService
}

func (ph *ProxyHandler) SaveSecretKEY(engine, key, apikey string) {
	result, err := ph.Service.SaveSecretKEY(engine, key, apikey)
	if err != nil {
		fmt.Println(result)
	}
	fmt.Println(result)
}

func (ph *ProxyHandler) ProxyGateway(endpoints domain.ProxyEndpoint, key, securityType string) {
	for _, endpoint := range endpoints.Endpoints {

		target, err := url.Parse(
			fmt.Sprintf("%s%s", endpoints.HostURI, endpoint.PathEndpoint),
		)
		if err != nil {
			log.Fatal(err)
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
					checkAPIKEYSecretKeyFromRequest(req, ph, key)
				}

			}
			proxy.ErrorHandler = func(rw http.ResponseWriter, r *http.Request, err error) {
				rw.WriteHeader(http.StatusInternalServerError)
				rw.Write([]byte(err.Error()))
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

func checkJWTSecretKeyFromRequest(req *http.Request, key string) error {
	header := req.Header.Get("Authorization")
	hs := jwt.NewHS256([]byte(key))
	now := time.Now()
	if !strings.HasPrefix(header, "Bearer ") {
		custonError := errors.NewError("Format is Authorization: Bearer [token]")
		logger.LogError(custonError.Error())
		return custonError
	}
	token := strings.Split(header, " ")[1]
	pl := JWTPayload{}
	expValidator := jwt.ExpirationTimeValidator(now)
	validatePayload := jwt.ValidatePayload(&pl.Payload, expValidator)
	_, err := jwt.Verify([]byte(token), hs, &pl, validatePayload)

	if errors.ErrorIs(err, jwt.ErrExpValidation) {
		logger.LogError(err.Error())
		return err
	}
	if errors.ErrorIs(err, jwt.ErrHMACVerification) {
		logger.LogError(err.Error())
		return err
	}

	return nil
}
func checkAPIKEYSecretKeyFromRequest(req *http.Request, ph *ProxyHandler, key string) {
	apikey, err := ph.Service.GetKEY(key)
	header := req.Header.Get("X-API-KEY")
	if err != nil {
		logger.LogError("getKey: failed " + err.Error())
	}
	if apikey == header {
		fmt.Println("OK")
	} else {
		fmt.Println("Invalid apikey")
	}
}

func modifyResponse(err error) func(*http.Response) error {
	return func(resp *http.Response) error {
		resp.Header.Set("X-Proxy", "EgoProxy")

		if err != nil {
			return err
		}
		return nil
	}
}
