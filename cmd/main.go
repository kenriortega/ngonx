package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/kenriortega/goproxy/cmd/cli"
	"github.com/kenriortega/goproxy/pkg/config"
	"github.com/kenriortega/goproxy/pkg/logger"
)

var (
	service        = "proxy"
	buildTime      string
	version        string
	versionHash    string
	configFromYaml config.Config
	errConfig      error
	portProxy      int
	generateApiKey = false
	serverList     = ""
	portLB         = 3030
	setingFile     = "goproxy.yaml"
	prevKey        = ""
	displayVersion = false
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
	flag.BoolVar(&displayVersion, "version", displayVersion, "Display version and exit")
	flag.StringVar(&service, "type", service, "Main Service default is proxy")
	flag.IntVar(&portProxy, "portProxy", portProxy, "Port to serve to run proxy")
	flag.BoolVar(&generateApiKey, "genkey", generateApiKey, "Action for generate hash for protected routes")
	flag.StringVar(&prevKey, "prevkey", prevKey, "Action for save a previous hash for protected routes to validate JWT")
	flag.StringVar(&serverList, "backends", serverList, "Load balanced backends, use commas to separate")
	flag.IntVar(&portLB, "portLB", portLB, "Port to serve to run load balancing")
	flag.Parse()

	if displayVersion {
		fmt.Printf("Version:\t%s\n", version)
		fmt.Printf("Version Git Hash:\t%s\n", versionHash)
		fmt.Printf("Build time:\t%s\n", buildTime)
		os.Exit(0)
	}

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
