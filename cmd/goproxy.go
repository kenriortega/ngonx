package main

import (
	"flag"
	"runtime"

	domain "github.com/kenriortega/goproxy/proxy/domain"

	"github.com/kenriortega/goproxy/cmd/cli"
	"github.com/kenriortega/goproxy/internal/infra"
	"github.com/kenriortega/goproxy/internal/utils"
)

var (
	service        = "proxy"
	config         infra.Config
	errConfig      error
	endpoints      []domain.ProxyEndpoint
	portProxy      int
	host           string
	generateApiKey = false
	serverList     = ""
	portLB         = 3030
	setingFile     = "goproxy.yaml"
	engine         = "badger"
	securityType   = "none"
)

func init() {
	config, errConfig = infra.LoadConfig(".", setingFile)
	if errConfig != nil {
		utils.LogError(errConfig.Error())
		utils.LogInfo("config: Creating setting file by default")
		// create empty file yml
		utils.CreateSettingFile(setingFile)
	}
	endpoints = config.ProxyGateway.EnpointsProxy
	portProxy = config.ProxyGateway.Port
	host = config.ProxyGateway.Host
	engine = config.ProxyCache.Engine
	securityType = config.ProxySecurity.Type
	generateApiKey = false

	numcpu := runtime.NumCPU()
	runtime.GOMAXPROCS(numcpu)
}

func main() {
	flag.StringVar(&service, "type", service, "Main Service default is proxy")
	flag.IntVar(&portProxy, "portProxy", portProxy, "Port to serve to run proxy")
	flag.BoolVar(&generateApiKey, "genkey", generateApiKey, "Action for generate hash for protected routes")
	flag.StringVar(&serverList, "backends", serverList, "Load balanced backends, use commas to separate")
	flag.IntVar(&portLB, "portLB", portLB, "Port to serve to run load balancing")
	flag.Parse()

	switch service {
	case "lb":
		cli.StartLB(serverList, portLB)
	case "proxy":
		cli.Start(generateApiKey, endpoints, host, portProxy, engine, securityType)
	}

}
