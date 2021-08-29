package cli

import (
	"net/http"

	"github.com/gorilla/mux"
	domain "github.com/kenriortega/ngonx/internal/mngt/domain"
	handlers "github.com/kenriortega/ngonx/internal/mngt/handlers"
	services "github.com/kenriortega/ngonx/internal/mngt/services"
	"github.com/kenriortega/ngonx/pkg/config"
	"github.com/kenriortega/ngonx/pkg/httpsrv"
)

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
