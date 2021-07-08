package cli

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/kenriortega/goproxy/internal/platform/badgerdb"
	"github.com/kenriortega/goproxy/internal/platform/logger"
	"github.com/kenriortega/goproxy/internal/platform/utils"
	domain "github.com/kenriortega/goproxy/internal/proxy/domain"
	handlers "github.com/kenriortega/goproxy/internal/proxy/handlers"
	services "github.com/kenriortega/goproxy/internal/proxy/services"
)

func Start(generateApiKey bool, endpoints []domain.ProxyEndpoint, host string, port int, engine, key, securityType string) {

	var proxyRepository domain.ProxyRepository
	clientBadger := badgerdb.GetBadgerDB(false)
	proxyRepository = domain.NewProxyRepository(clientBadger)
	h := handlers.ProxyHandler{
		Service: services.NewProxyService(proxyRepository),
	}

	if generateApiKey {
		word := utils.StringWithCharset()
		apiKey := utils.ApiKeyGenerator(word)
		_, err := h.Service.SaveSecretKEY(engine, key, apiKey)
		if err != nil {
			logger.LogError("genkey: Failed " + err.Error())
		}
		logger.LogInfo("genkey: Susscefull")
	}

	for _, endpoints := range endpoints {

		h.ProxyGateway(endpoints, key, securityType)
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
