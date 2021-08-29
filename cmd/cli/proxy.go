package cli

import (
	domain "github.com/kenriortega/ngonx/internal/proxy/domain"
	handlers "github.com/kenriortega/ngonx/internal/proxy/handlers"
	services "github.com/kenriortega/ngonx/internal/proxy/services"
	"github.com/kenriortega/ngonx/pkg/badgerdb"
	"github.com/kenriortega/ngonx/pkg/config"
	"github.com/kenriortega/ngonx/pkg/genkey"
	"github.com/kenriortega/ngonx/pkg/httpsrv"
	"github.com/kenriortega/ngonx/pkg/logger"
)

func StartProxy(
	generateApiKey bool,
	port int,
	prevKey string,
	config config.Config,
) {

	engine := config.ProxyCache.Engine
	securityType := config.ProxySecurity.Type
	key := config.ProxyCache.Key + "_" + securityType

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

	for _, endpoints := range config.ProxyGateway.EnpointsProxy {

		h.ProxyGateway(endpoints, engine, key, securityType)
	}

	if config.ProxySSL.Enable {
		portSSL := config.ProxyGateway.Port + config.ProxySSL.SSLPort
		server := httpsrv.NewServerSSL(
			config.ProxyGateway.Host,
			portSSL,
			nil,
		)
		server.StartSSL(
			config.ProxySSL.CrtFile,
			config.ProxySSL.KeyFile,
		)
	} else {
		port = config.ProxyGateway.Port + port
		server := httpsrv.NewServer(
			config.ProxyGateway.Host,
			port,
			nil,
		)
		server.Start()
	}

}
