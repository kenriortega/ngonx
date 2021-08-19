package metric

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	CountersByRoute = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "counter_by_routes",
			Help: "Register all access by routes using counter metrics",
		},
		[]string{
			"service",
			"path",
		},
	)
	CountersByEndpoint = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "counter_request_by_microservicio",
			Help: "Register all call to the endpoints",
		},
		[]string{
			"proxyPath",
			"endpointPath",
			"ipAddr",
			"method",
		},
	)

	DurationHttpRequest = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_response_time_seconds",
		Help: "Duration of HTTP requests.",
	}, []string{"path"})

	TotalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Number of get requests.",
		},
		[]string{"path"},
	)

	ResponseStatus = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "response_status",
			Help: "Status of HTTP response",
		},
		[]string{"status"},
	)
)

func init() {
	prometheus.MustRegister(CountersByRoute)
	prometheus.MustRegister(CountersByEndpoint)
	prometheus.MustRegister(DurationHttpRequest)
	prometheus.MustRegister(TotalRequests)
	prometheus.MustRegister(ResponseStatus)
}

func ExposeMetricServer(configPort int) {
	http.Handle("/metrics", promhttp.Handler())
	port := fmt.Sprintf(":%d", configPort)
	log.Fatal(http.ListenAndServe(port, nil))
}
