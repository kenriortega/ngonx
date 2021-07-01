package proxy

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	domain "egosystem.org/micros/gateway/domain"
	handlers "egosystem.org/micros/gateway/handlers"
	"egosystem.org/micros/internal"
)

var (
	config    internal.Config
	errConfig error
	endpoints []domain.EndpointService
	port      int
	host      string
)

func init() {
	config, errConfig = internal.LoadConfig(".", "proxy.yaml")
	if errConfig != nil {
		log.Println(errConfig)
	}
	endpoints = config.ProxyGateway.EnpointsProxy[0].Services
	port = config.ProxyGateway.Port
	host = config.ProxyGateway.Host
}

func Start() {
	flag.IntVar(&port, "port", port, "Port to serve")
	flag.Parse()

	for _, endpoint := range endpoints {
		handlers.ProxyGateway(endpoint)
	}

	server := &http.Server{
		Handler: nil,
		Addr:    fmt.Sprintf("%s:%d", host, port),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Load Balancer started at :%d\n", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
