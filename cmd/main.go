package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/kenriortega/ngonx/cmd/cli"
	"github.com/kenriortega/ngonx/pkg/config"
	"github.com/kenriortega/ngonx/pkg/logger"
	"github.com/kenriortega/ngonx/pkg/metric"
)

var (
	buildTime           string
	version             string
	versionHash         string
	configFromYaml      config.Config
	errConfig           error
	portProxy           int
	service             = "proxy"
	setup               = false
	generateApiKey      = false
	serverList          = ""
	portLB              = 3030
	pathSettingsFile, _ = os.Getwd()
	settingsFile        = "ngonx.yaml"
	prevKey             = ""
	displayVersion      = false
	portExporterProxy   = 10000
)

func init() {

	flag.BoolVar(&displayVersion, "version", displayVersion, "Display version and exit")
	flag.BoolVar(&setup, "setup", setup, "Create yaml file configuration")
	flag.StringVar(&pathSettingsFile, "configPath", pathSettingsFile, "Config path only not filename.yaml")
	flag.StringVar(&settingsFile, "configFile", settingsFile, "Only config filename.yaml default ngonx.yaml")
	flag.StringVar(&service, "type", service, "Main Service default is proxy")
	flag.IntVar(&portProxy, "proxyPort", portProxy, "Port to serve to run proxy")
	flag.BoolVar(&generateApiKey, "genkey", generateApiKey, "Action for generate hash for protected routes")
	flag.StringVar(&prevKey, "prevkey", prevKey, "Action for save a previous hash for protected routes to validate JWT")
	flag.StringVar(&serverList, "backends", serverList, "Load balanced backends, use commas to separate")
	flag.IntVar(&portLB, "lbPort", portLB, "Port to serve to run load balancing default 3030")
	flag.IntVar(&portExporterProxy, "portExporter", portExporterProxy, "Port to serve expose metrics from prometheus default 10000")
	flag.Parse()

	numcpu := runtime.NumCPU()
	runtime.GOMAXPROCS(numcpu)
}

func main() {

	if setup {
		logger.LogInfo("config: Creating setting file by default")
		configFromYaml.CreateSettingFile(settingsFile)
		os.Exit(0)
	} else {
		configFromYaml, errConfig = config.LoadConfig(pathSettingsFile, settingsFile)

		if errConfig != nil {
			logger.LogError("Yaml file not found please run command setup " + errConfig.Error())
		}
	}

	if displayVersion {
		fmt.Printf("Version:\t%s\n", version)
		fmt.Printf("Version Git Hash:\t%s\n", versionHash)
		fmt.Printf("Build time:\t%s\n", buildTime)
		os.Exit(0)
	}
	// Exporter Metrics
	go metric.ExposeMetricServer(portExporterProxy)
	// Admin pannel
	go cli.StartMngt(configFromYaml)

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
