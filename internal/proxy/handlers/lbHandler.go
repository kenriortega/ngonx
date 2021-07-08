package proxy

import (
	"log"
	"net/http"
	"time"

	domain "github.com/kenriortega/goproxy/internal/proxy/domain"
)

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
		log.Printf("%s(%s) Max attempts reached, terminating\n", r.RemoteAddr, r.URL.Path)
		http.Error(w, "Service not available", http.StatusServiceUnavailable)
		return
	}

	peer := ServerPool.GetNextPeer()
	if peer != nil {
		peer.ReverseProxy.ServeHTTP(w, r)
		return
	}
	http.Error(w, "Service not available", http.StatusServiceUnavailable)
}

// healthCheck runs a routine for check status of the backends every 2 mins
func HealthCheck() {
	t := time.NewTicker(time.Minute * 1)
	for range t.C {
		log.Println("Starting health check...")
		ServerPool.HealthCheck()
		log.Println("Health check completed")
	}
}
