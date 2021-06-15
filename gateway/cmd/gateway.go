package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type EndpointService struct {
	HostURI string `json:"host_url"`
	Path    string `json:"path"`
}

func main() {
	// TODO: migrate to yaml | json | key/value service
	endpoints := []EndpointService{
		{
			HostURI: "http://localhost:8000/api/v1/health/",
			Path:    "/health",
		},
		{
			HostURI: "http://localhost:8000/api/v1/version/",
			Path:    "/version",
		},
	}

	for _, endpoint := range endpoints {
		ManagementAPI(endpoint.HostURI, endpoint.Path)
	}
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func ManagementAPI(endpoint, path string) {
	target, err := url.Parse(endpoint)
	if err != nil {
		log.Fatal(err)
	}
	http.Handle(
		path,
		httputil.NewSingleHostReverseProxy(target),
	)
}
