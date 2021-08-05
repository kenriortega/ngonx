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

	"github.com/kenriortega/goproxy/pkg/logger"

	domain "github.com/kenriortega/goproxy/internal/proxy/domain"

	handlers "github.com/kenriortega/goproxy/internal/proxy/handlers"
)

func StartLB(serverList string, port int) {

	if len(serverList) == 0 {
		log.Fatal("Please provide one or more backends to load balance")
	}

	// parse servers
	tokens := strings.Split(serverList, ",")
	for _, tok := range tokens {
		serverUrl, err := url.Parse(tok)
		if err != nil {
			logger.LogError(err.Error())
		}

		proxy := httputil.NewSingleHostReverseProxy(serverUrl)
		proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, e error) {
			logger.LogInfo(fmt.Sprintf("[%s] %s\n", serverUrl.Host, e.Error()))
			retries := handlers.GetRetryFromContext(request)
			if retries < 3 {
				// It looks like this should be a
				time.Sleep(10 * time.Millisecond)
				ctx := context.WithValue(request.Context(), domain.RETRY, retries+1)
				proxy.ServeHTTP(writer, request.WithContext(ctx))
				return
			}

			// after 3 retries, mark this backend as down
			handlers.ServerPool.MarkBackendStatus(serverUrl, false)

			// if the same request routing for few attempts with different backends, increase the count
			attempts := handlers.GetAttemptsFromContext(request)
			logger.LogInfo(fmt.Sprintf("%s(%s) Attempting retry %d\n", request.RemoteAddr, request.URL.Path, attempts))
			ctx := context.WithValue(request.Context(), domain.ATTEMPTS, attempts+1)
			handlers.Lbalancer(writer, request.WithContext(ctx))
		}

		handlers.ServerPool.AddBackend(&domain.Backend{
			URL:          serverUrl,
			Alive:        true,
			ReverseProxy: proxy,
		})
		logger.LogInfo(fmt.Sprintf("Configured server: %s\n", serverUrl))
	}

	// create http server
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: http.HandlerFunc(handlers.Lbalancer),
	}

	// start health checking
	go handlers.HealthCheck()

	logger.LogInfo(fmt.Sprintf("Load Balancer started at :%d\n", port))
	if err := server.ListenAndServe(); err != nil {
		logger.LogError(err.Error())
	}
}
