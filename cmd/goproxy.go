package main

import (
	"flag"
	"log"
	"runtime"

	domain "github.com/kenriortega/goproxy/proxy/domain"

	"github.com/kenriortega/goproxy/cmd/cli"
	"github.com/kenriortega/goproxy/internal/infra"
)

var ()

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
)

func init() {
	config, errConfig = infra.LoadConfig(".", "goproxy.yaml")
	if errConfig != nil {
		log.Println(errConfig)
	}
	endpoints = config.ProxyGateway.EnpointsProxy
	portProxy = config.ProxyGateway.Port
	host = config.ProxyGateway.Host
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
		cli.Start(generateApiKey, endpoints, host, portProxy)
	}

}
