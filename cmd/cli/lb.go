package cli

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	domain "github.com/kenriortega/ngonx/internal/proxy/domain"

	handlers "github.com/kenriortega/ngonx/internal/proxy/handlers"
	"github.com/kenriortega/ngonx/pkg/backoff"
	"github.com/kenriortega/ngonx/pkg/errors"
	"github.com/kenriortega/ngonx/pkg/logger"
	"github.com/spf13/cobra"
)

// MaxJitter will randomize over the full exponential backoff time
const MaxJitter = 1.0

// NoJitter disables the use of jitter for randomizing the exponential backoff time
const NoJitter = 0.0

var lbCmd = &cobra.Command{
	Use:   "lb",
	Short: "Run ngonx as a load balancer (round robin)",
	Run: func(cmd *cobra.Command, args []string) {
		port, err := cmd.Flags().GetInt(flagPort)
		if err != nil {
			logger.LogError(errors.Errorf("lb: %v", err).Error())
		}
		serverList, err := cmd.Flags().GetString(flagServerList)
		if err != nil {
			log.Fatalf(err.Error())
		}
		if len(serverList) == 0 {
			logger.LogError(errors.Errorf("lb: provide one or more backends to load balance %v", err).Error())
		}

		// parse servers
		tokens := strings.Split(serverList, ",")
		for _, tok := range tokens {
			serverUrl, err := url.Parse(tok)
			if err != nil {
				logger.LogError(errors.Errorf("lb: %v", err).Error())
			}

			proxy := httputil.NewSingleHostReverseProxy(serverUrl)
			proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, e error) {
				logger.LogInfo(fmt.Sprintf("lb: %s %s\n", serverUrl.Host, e.Error()))
				retry := handlers.GetRetryFromContext(request)

				if retry < 3 {
					time.Sleep(backoff.Default.Duration(retry))
					ctx := context.WithValue(request.Context(), domain.RETRY, retry+1)
					proxy.ServeHTTP(writer, request.WithContext(ctx))

					return
				}

				// after 3 retries, mark this backend as down
				handlers.ServerPool.MarkBackendStatus(serverUrl, false)

				// if the same request routing for few attempts with different backends, increase the count
				attempts := handlers.GetAttemptsFromContext(request)
				logger.LogInfo(fmt.Sprintf("lb: %s(%s) Attempting retry %d\n", request.RemoteAddr, request.URL.Path, attempts))
				ctx := context.WithValue(request.Context(), domain.ATTEMPTS, attempts+1)
				handlers.Lbalancer(writer, request.WithContext(ctx))
			}

			handlers.ServerPool.AddBackend(&domain.Backend{
				URL:          serverUrl,
				Alive:        true,
				ReverseProxy: proxy,
			})
			logger.LogInfo(fmt.Sprintf("lb: configured server: %s\n", serverUrl))
		}

		// create http server
		server := http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: http.HandlerFunc(handlers.Lbalancer),
		}

		// start health checking
		go handlers.HealthCheck()

		logger.LogInfo(fmt.Sprintf("lb: Load Balancer started at :%d\n", port))
		if err := server.ListenAndServe(); err != nil {
			logger.LogError(errors.Errorf("lb: %v", err).Error())
		}

	},
}

func init() {
	lbCmd.Flags().String(flagServerList, cfgFile, "Load balanced backends, use commas to separate")
	lbCmd.Flags().Int(flagPort, 4000, "Port to serve to run load balancing ")

	rootCmd.AddCommand(lbCmd)
}
