package gateway

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	domain "egosystem.org/micros/gateway/domain"
	services "egosystem.org/micros/gateway/services"
	"egosystem.org/micros/internal"
)

type ProxyHandler struct {
	Service services.DefaultProxyService
}

func (ph *ProxyHandler) SaveSecretKEY(engine, apikey string) {
	result, err := ph.Service.SaveSecretKEY(engine, apikey)
	if err != nil {
		fmt.Println(result)
	}
	fmt.Println(result)
}

func ProxyGateway(endpoints domain.ProxyEndpoint) {
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
			fmt.Println(req)
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

// GetNameRandom get string names
func GenRandomKEY(length int) string {
	word := internal.StringWithCharset(length, internal.CHARSET)
	apiKey := internal.ApiKeyGenerator(word)
	return apiKey
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
