package cli

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	domain "github.com/kenriortega/ngonx/internal/mngt/domain"
	handlers "github.com/kenriortega/ngonx/internal/mngt/handlers"
	services "github.com/kenriortega/ngonx/internal/mngt/services"
	"github.com/kenriortega/ngonx/pkg/config"
	"github.com/kenriortega/ngonx/pkg/httpsrv"
	"github.com/kenriortega/ngonx/pkg/logger"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ngonxctl",
	Short: "A proxy reverse inspired on nginx & traefik",
	Long:  `This is Ngonx ctl a proxy reverse inspired on nginx & traefik`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {

		logger.LogError(err.Error())
		os.Exit(1)
	}
	go StartMngt(configFromYaml)

}
func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, flagCfgFile, "f", cfgFile, "File setting.yml")
	rootCmd.PersistentFlags().StringVarP(&cfgPath, flagCfgPath, "p", cfgPath, "Config path only ")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	configFromYaml, errConfig = config.LoadConfig(cfgPath, cfgFile)

	if errConfig != nil {
		logger.LogError("Yaml file not found please run command setup " + errConfig.Error())
	}
}

func StartMngt(config config.Config) {
	r := mux.NewRouter()
	repo := domain.NewMngtRepositoryStorage()
	service := services.NewMngtService(repo)
	mh := handlers.NewMngtHandler(service)

	// populate data from config file list of services

	for _, endpoints := range config.ProxyGateway.EnpointsProxy {
		for _, endpoint := range endpoints.Endpoints {
			endpointMap := make(map[string]interface{})
			endpointMap["path_url"] = endpoint.PathEndpoint
			endpointMap["status"] = "down"
			mh.RegisterEnpoint(endpointMap)
		}
	}
	// Routes...
	adminRoutes := r.PathPrefix("/api/v1/mngt").Subrouter()
	adminRoutes.HandleFunc("/", mh.GetAllEndpoints).Methods(http.MethodGet)
	port := 10_001
	server := httpsrv.NewServer(
		"0.0.0.0",
		port,
		r,
	)
	server.Start()
}
