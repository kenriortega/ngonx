package proxy

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	domain "github.com/kenriortega/goproxy/proxy/domain"
	services "github.com/kenriortega/goproxy/proxy/services"
)

var proxy *httputil.ReverseProxy

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

func (ph *ProxyHandler) ProxyGateway(endpoints domain.ProxyEndpoint, securityType string) {
	for _, endpoint := range endpoints.Endpoints {

		target, err := url.Parse(
			fmt.Sprintf("%s%s", endpoints.HostURI, endpoint.PathEndpoint),
		)
		if err != nil {
			log.Fatal(err)
		}
		if endpoint.PathProtected {
			proxy = httputil.NewSingleHostReverseProxy(target)
			proxy.ModifyResponse = modifyResponse()

			originalDirector := proxy.Director
			proxy.Director = func(req *http.Request) {
				originalDirector(req)

				switch securityType {
				case "jwt":
					checkJWTSecretKeyFromRequest(req)
				case "apikey":
					checkAPIKEYSecretKeyFromRequest(req)
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
			proxy.ModifyResponse = modifyResponse()

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

func checkJWTSecretKeyFromRequest(req *http.Request) {
	header := req.Header.Get("Authorization")
	fmt.Println(header)
	// secretKey, err := ph.Service.GetKEY("secretKey")
	// if err != nil {
	// 	utils.LogError("getKey: failed " + err.Error())
	// }
	// fmt.Println(secretKey)
}
func checkAPIKEYSecretKeyFromRequest(req *http.Request) {
	header := req.Header.Get("X-API-KEY")
	fmt.Println(header)
}

func modifyResponse() func(*http.Response) error {
	return func(resp *http.Response) error {
		resp.Header.Set("X-Proxy", "EgoProxy")
		return nil
	}
}
