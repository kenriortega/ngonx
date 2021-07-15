package cli

import (
	"fmt"
	"net/http"
	"time"

	"github.com/kenriortega/goproxy/internal/platform/badgerdb"
	"github.com/kenriortega/goproxy/internal/platform/genkey"
	"github.com/kenriortega/goproxy/internal/platform/logger"
	domain "github.com/kenriortega/goproxy/internal/proxy/domain"
	handlers "github.com/kenriortega/goproxy/internal/proxy/handlers"
	services "github.com/kenriortega/goproxy/internal/proxy/services"
)

func Start(generateApiKey bool, endpoints []domain.ProxyEndpoint, host string, port int, engine, key, prevKey, securityType string) {

	var proxyRepository domain.ProxyRepository
	clientBadger := badgerdb.GetBadgerDB(false)
	proxyRepository = domain.NewProxyRepository(clientBadger)
	h := handlers.ProxyHandler{
		Service: services.NewProxyService(proxyRepository),
	}

	if generateApiKey {
		word := genkey.StringWithCharset()
		apiKey := genkey.ApiKeyGenerator(word)
		_, err := h.Service.SaveSecretKEY(engine, key, apiKey)
		if err != nil {
			logger.LogError("genkey: Failed " + err.Error())
		}
		logger.LogInfo("genkey: Susscefull")
	}
	if prevKey != "" {
		_, err := h.Service.SaveSecretKEY(engine, key, prevKey)
		if err != nil {
			logger.LogError("prevKey: Failed " + err.Error())
		}
		logger.LogInfo("prevKey: Susscefull")
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

	logger.LogInfo(fmt.Sprintf("Proxy started at :%d\n", port))
	if err := server.ListenAndServe(); err != nil {
		logger.LogError(err.Error())
	}
}
