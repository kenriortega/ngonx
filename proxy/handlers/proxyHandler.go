package proxy

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/kenriortega/goproxy/internal/utils"
	domain "github.com/kenriortega/goproxy/proxy/domain"
	services "github.com/kenriortega/goproxy/proxy/services"
)

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

func (ph *ProxyHandler) ProxyGateway(endpoints domain.ProxyEndpoint) {
	for _, endpoint := range endpoints.Endpoints {

		target, err := url.Parse(
			fmt.Sprintf("%s%s", endpoints.HostURI, endpoint.PathEndpoint),
		)
		if err != nil {
			log.Fatal(err)
		}
		proxy := httputil.NewSingleHostReverseProxy(target)
		proxy.ModifyResponse = modifyResponse()

		originalDirector := proxy.Director
		proxy.Director = func(req *http.Request) {
			originalDirector(req)
			secretKey, err := ph.Service.GetKEY("secretKey")
			if err != nil {
				utils.LogError("getKey: failed " + err.Error())
			}
			fmt.Println(secretKey)
			modifyRequest(req)
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

func modifyRequest(req *http.Request) {
	req.Header.Set("X-Proxy", "Simple-Reverse-Proxy")
}
func modifyResponse() func(*http.Response) error {
	return func(resp *http.Response) error {
		resp.Header.Set("X-Proxy", "EgoProxy")
		return nil
	}
}
