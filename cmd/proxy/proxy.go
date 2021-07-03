package proxy

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	domain "egosystem.org/micros/gateway/domain"
	handlers "egosystem.org/micros/gateway/handlers"
	services "egosystem.org/micros/gateway/services"
	"egosystem.org/micros/internal"
)

var (
	config         Config
	errConfig      error
	endpoints      []domain.ProxyEndpoint
	port           int
	host           string
	generateApiKey bool
)

func init() {
	config, errConfig = LoadConfig(".", "proxy.yaml")
	if errConfig != nil {
		log.Println(errConfig)
	}
	endpoints = config.ProxyGateway.EnpointsProxy
	port = config.ProxyGateway.Port
	host = config.ProxyGateway.Host
	generateApiKey = false
}

func Start() {
	flag.IntVar(&port, "port", port, "Port to serve")
	flag.BoolVar(&generateApiKey, "genkey", generateApiKey, "Action for generate hash")
	flag.Parse()

	var proxyRepository domain.ProxyRepository
	clientBadger := GetBadgerDB(true)
	proxyRepository = domain.NewProxyRepository(clientBadger)
	h := handlers.ProxyHandler{
		Service: services.NewProxyService(proxyRepository),
	}

	if generateApiKey {
		word := internal.StringWithCharset()
		apiKey := internal.ApiKeyGenerator(word)
		_, err := h.Service.SaveSecretKEY("badger", apiKey)
		if err != nil {
			internal.LogError("genkey: Failed " + err.Error())
		}
		internal.LogInfo("genkey: Susscefull")
	}

	for _, endpoints := range endpoints {

		h.ProxyGateway(endpoints)
	}

	server := &http.Server{
		Handler: nil,
		Addr:    fmt.Sprintf("%s:%d", host, port),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Proxy started at :%d\n", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
