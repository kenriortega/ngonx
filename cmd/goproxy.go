package main

import (
	"flag"
	"runtime"

	"github.com/kenriortega/goproxy/cmd/cli"
	"github.com/kenriortega/goproxy/internal/platform/config"
	"github.com/kenriortega/goproxy/internal/platform/logger"
)

var (
	service        = "proxy"
	configFromYaml config.Config
	errConfig      error
	portProxy      int
	generateApiKey = false
	serverList     = ""
	portLB         = 3030
	setingFile     = "goproxy.yaml"
	prevKey        = ""
)

func init() {
	configFromYaml, errConfig = config.LoadConfig(".", setingFile)
	if errConfig != nil {
		logger.LogError(errConfig.Error())
		logger.LogInfo("config: Creating setting file by default")
		// create empty file yml
		configFromYaml.CreateSettingFile(setingFile)
	}

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
		cli.StartProxy(
			generateApiKey,
			portProxy,
			prevKey,
			configFromYaml,
		)
	case "static":
		cli.StartStaticServer(configFromYaml)
	}

}
