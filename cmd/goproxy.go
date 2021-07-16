package main

import (
	"flag"
	"runtime"

	"github.com/kenriortega/goproxy/cmd/cli"
	"github.com/kenriortega/goproxy/internal/platform/config"
	"github.com/kenriortega/goproxy/internal/platform/logger"
	domain "github.com/kenriortega/goproxy/internal/proxy/domain"
)

var (
	service        = "proxy"
	configFromYaml config.Config
	errConfig      error
	endpoints      []domain.ProxyEndpoint
	portProxy      int
	host           string
	generateApiKey = false
	serverList     = ""
	portLB         = 3030
	setingFile     = "goproxy.yaml"
	engine         = "badger"
	key            = "secretKey"
	securityType   = "none"
	prevKey        = ""
	sslProxy       config.ProxySSL
)

func init() {
	configFromYaml, errConfig = config.LoadConfig(".", setingFile)
	if errConfig != nil {
		logger.LogError(errConfig.Error())
		logger.LogInfo("config: Creating setting file by default")
		// create empty file yml
		configFromYaml.CreateSettingFile(setingFile)
	}
	endpoints = configFromYaml.ProxyGateway.EnpointsProxy
	portProxy = configFromYaml.ProxyGateway.Port
	host = configFromYaml.ProxyGateway.Host
	engine = configFromYaml.ProxyCache.Engine
	securityType = configFromYaml.ProxySecurity.Type
	key = configFromYaml.ProxyCache.Key + "_" + securityType
	sslProxy = configFromYaml.ProxySSL
	generateApiKey = false

	numcpu := runtime.NumCPU()
	runtime.GOMAXPROCS(numcpu)
}

func main() {
	flag.StringVar(&service, "type", service, "Main Service default is proxy")
	flag.IntVar(&portProxy, "portProxy", portProxy, "Port to serve to run proxy")
	flag.BoolVar(&generateApiKey, "genkey", generateApiKey, "Action for generate hash for protected routes")
	flag.StringVar(&prevKey, "prevkey", prevKey, "Action for save a previous hash for protected routes to validate JWT")
	flag.StringVar(&serverList, "backends", serverList, "Load balanced backends, use commas to separate")
	flag.IntVar(&portLB, "portLB", portLB, "Port to serve to run load balancing")
	flag.Parse()

	switch service {
	case "lb":
		cli.StartLB(serverList, portLB)
	case "proxy":
		cli.Start(
			generateApiKey,
			endpoints,
			host,
			portProxy,
			engine,
			key,
			prevKey,
			securityType,
			sslProxy,
		)
	}

}
