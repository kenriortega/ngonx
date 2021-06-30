package proxy

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	domain "egosystem.org/micros/gateway/domain"
	handler "egosystem.org/micros/gateway/handlers"
)

func Start() {
	var port int
	flag.IntVar(&port, "port", 8081, "Port to serve")
	flag.Parse()
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
		Addr:    fmt.Sprintf("0.0.0.0:%d", port),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
