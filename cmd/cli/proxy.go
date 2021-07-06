package cli

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/kenriortega/goproxy/internal/infra"
	"github.com/kenriortega/goproxy/internal/utils"
	domain "github.com/kenriortega/goproxy/proxy/domain"
	handlers "github.com/kenriortega/goproxy/proxy/handlers"
	services "github.com/kenriortega/goproxy/proxy/services"
)

var (
	config         infra.Config
	errConfig      error
	endpoints      []domain.ProxyEndpoint
	port           int
	host           string
	generateApiKey bool
)

func init() {
	config, errConfig = infra.LoadConfig(".", "proxy.yaml")
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
	clientBadger := infra.GetBadgerDB(false)
	proxyRepository = domain.NewProxyRepository(clientBadger)
	h := handlers.ProxyHandler{
		Service: services.NewProxyService(proxyRepository),
	}

	if generateApiKey {
		word := utils.StringWithCharset()
		apiKey := utils.ApiKeyGenerator(word)
		_, err := h.Service.SaveSecretKEY("badger", "secretKey", apiKey)
		if err != nil {
			utils.LogError("genkey: Failed " + err.Error())
		}
		utils.LogInfo("genkey: Susscefull")
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
