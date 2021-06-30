package gateway

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	domain "egosystem.org/micros/gateway/domain"
)

func ProxyGateway(endpoint domain.EndpointService) {
	target, err := url.Parse(endpoint.HostURI)
	if err != nil {
		log.Fatal(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(target)
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		modifyRequest(req)
	}

	proxy.ErrorHandler = errorHandler()
	http.Handle(
		endpoint.Path,
		http.StripPrefix(
			endpoint.Path,
			proxy,
		),
	)
}

func modifyRequest(req *http.Request) {
	req.Header.Set("X-Proxy", "Egosystem-Proxy")
}

func errorHandler() func(http.ResponseWriter, *http.Request, error) {
	return func(w http.ResponseWriter, req *http.Request, err error) {
		fmt.Printf("Got error while modifying response: %v \n", err)
	}
}
