package proxy

import (
	"fmt"
	"net/http"
	"time"

	domain "github.com/kenriortega/goproxy/internal/proxy/domain"
	"github.com/kenriortega/goproxy/pkg/errors"
	"github.com/kenriortega/goproxy/pkg/logger"
)

// ServerPool struct for server pool
var ServerPool domain.ServerPool

// GetAttemptsFromContext returns the attempts for request
func GetAttemptsFromContext(r *http.Request) int {
	if attempts, ok := r.Context().Value(domain.ATTEMPTS).(int); ok {
		return attempts
	}
	return 1
}

// GetAttemptsFromContext returns the attempts for request
func GetRetryFromContext(r *http.Request) int {
	if retry, ok := r.Context().Value(domain.RETRY).(int); ok {
		return retry
	}
	return 0
}

// Lbalancer load balances the incoming request
func Lbalancer(w http.ResponseWriter, r *http.Request) {
	attempts := GetAttemptsFromContext(r)
	if attempts > 3 {
		logger.LogInfo(fmt.Sprintf("%s(%s) Max attempts reached, terminating\n", r.RemoteAddr, r.URL.Path))
		http.Error(w, errors.ErrLBHttp.Error(), http.StatusServiceUnavailable)
		return
	}

	peer := ServerPool.GetNextPeer()
	if peer != nil {
		peer.ReverseProxy.ServeHTTP(w, r)
		return
	}
	http.Error(w, errors.ErrLBHttp.Error(), http.StatusServiceUnavailable)
}

// HealthCheck runs a routine for check status of the backends every 2 mins
func HealthCheck() {
	t := time.NewTicker(time.Minute * 1)
	for range t.C {
		logger.LogInfo("Starting health check...")
		ServerPool.HealthCheck()
		logger.LogInfo("Health check completed")
	}
}
