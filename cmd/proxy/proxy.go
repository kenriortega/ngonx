package proxy

import (
	"log"
	"net/http"
	"time"

	domain "egosystem.org/micros/gateway/domain"
	handler "egosystem.org/micros/gateway/handlers"
)

func Start() {
	// TODO: migrate to yaml | json | key/value service
	endpoints := []domain.EndpointService{
		{
			HostURI: "http://localhost:8000/api/v1/health/",
			Path:    "/health/",
		},
		{
			HostURI: "http://localhost:8000/api/v1/version/",
			Path:    "/version/",
		},
	}

	for _, endpoint := range endpoints {
		handler.ProxyGateway(endpoint)
	}

	srv := &http.Server{
		Handler: nil,
		Addr:    "0.0.0.0:8080",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
