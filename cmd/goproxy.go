package main

import (
	"flag"

	"github.com/kenriortega/goproxy/cmd/cli"
)

var (
	service = "proxy"
)

func main() {
	flag.StringVar(&service, "type", service, "Service type default is proxy mode")

	switch service {
	case "lb":
		cli.StartLB()
	default:
		cli.Start()
	}
}
