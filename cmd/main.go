package main

import (
	"runtime"

	"github.com/kenriortega/ngonx/cmd/cli"
)

func init() {

	numcpu := runtime.NumCPU()
	runtime.GOMAXPROCS(numcpu) // Try to use all available CPUs.
}

func main() {
	cli.Execute()
}
