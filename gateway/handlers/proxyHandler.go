package gateway

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	domain "egosystem.org/micros/gateway/domain"
)

func ProxyGateway(endpoint domain.EndpointService) {
	target, err := url.Parse(endpoint.HostURI)
	if err != nil {
		log.Fatal(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	proxy.Director = func(req *http.Request) {
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Set("X-Proxy", "Egosystem-Proxy")
		req.Header.Add("X-Origin-Host", target.Host)
	}

	proxy.ErrorHandler = func(rw http.ResponseWriter, r *http.Request, err error) {
		fmt.Printf("error was: %+v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
	}

	proxy.Transport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
	}

	proxy.ModifyResponse = func(r *http.Response) error {
		return nil
	}

	http.Handle(
		endpoint.Path,
		http.StripPrefix(
			endpoint.Path,
			proxy,
		),
	)
}
