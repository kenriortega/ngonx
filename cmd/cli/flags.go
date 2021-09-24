package cli

import (
	"os"

	"github.com/kenriortega/ngonx/pkg/config"
)

var (
	// variables to store data for global usage
	configFromYaml config.Config
	cfgFile        = "ngonx.yaml"
	cfgPath, _     = os.Getwd()
	errConfig      error

	// flags
	flagPort       = "port"
	flagServerList = "backends"
	flagGenApiKey  = "genkey"
	flagPrevKey    = "prevkey"
	flagCfgFile    = "cfgfile"
	flagCfgPath    = "cfgpath"
)
