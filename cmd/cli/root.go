package cli

import (
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/mux"
	domain "github.com/kenriortega/ngonx/internal/mngt/domain"
	handlers "github.com/kenriortega/ngonx/internal/mngt/handlers"
	services "github.com/kenriortega/ngonx/internal/mngt/services"
	"github.com/kenriortega/ngonx/pkg/config"
	"github.com/kenriortega/ngonx/pkg/healthcheck"
	"github.com/kenriortega/ngonx/pkg/httpsrv"
	"github.com/kenriortega/ngonx/pkg/logger"
	"github.com/spf13/cobra"
)

// Middleware CORS
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
		next.ServeHTTP(w, r)
	})
}

// SSE logic

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

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
	go StartMngt(configFromYaml)

}

func StartMngt(config config.Config) {
	r := mux.NewRouter()
	repo := domain.NewMngtRepositoryStorage()
	service := services.NewMngtService(repo)
	mh := handlers.NewMngtHandler(service)

	// populate data from config file list of services

	for _, endpoints := range config.ProxyGateway.EnpointsProxy {
		hostUri := endpoints.HostURI
		for _, it := range endpoints.Endpoints {
			endpointMap := make(map[string]interface{})
			endpointMap["path_url"] = hostUri + it.PathEndpoint
			endpointMap["status"] = "down"
			mh.RegisterEndpoint(endpointMap)
		}
	}
	// Routes...
	mngt := r.PathPrefix("/api/v1/mngt").Subrouter()
	mngt.HandleFunc("/", mh.GetAllEndpoints).Methods(http.MethodGet)
	mngt.HandleFunc("/health", healthHandler)
	mngt.HandleFunc("/readiness", readinessHandler)

	// Realtime options
	mngt.HandleFunc("/wss", mh.WssocketHandler)
	r.Use(CORSMiddleware)
	port := 10_001
	server := httpsrv.NewServer(
		"0.0.0.0",
		port,
		r,
	)

	go func() {
		t := time.NewTicker(time.Second * 30)
		for range t.C {

			endpoints, err := service.ListEndpoints()
			if err != nil {
				logger.LogError(err.Error())
			}
			for _, it := range endpoints {
				u, err := url.Parse(it.PathUrl)
				if err != nil {
					logger.LogError(err.Error())
				}
				status := healthcheck.IsBackendAlive(u)
				if status {
					it.Status = "up"
				} else {
					it.Status = "down"
				}
				mh.UpdateEndpoint(it)
			}

		}
	}()

	server.Start()
}
