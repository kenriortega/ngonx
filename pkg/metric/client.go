package metric

import (
	"fmt"
	"log"
	"net/http"

	"github.com/kenriortega/ngonx/pkg/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
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
		Name: "http_response_time_seconds_microservices",
		Help: "Duration of HTTP requests.",
	}, []string{"path"})

	TotalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Number of get requests.",
		},
		[]string{"path", "service"},
	)
)

func init() {
	err := prometheus.Register(DurationHttpRequest)
	if err != nil {
		logger.LogError(err.Error())
	}
	prometheus.MustRegister(CountersByEndpoint)
	prometheus.MustRegister(TotalRequests)
}

func ExposeMetricServer(configPort int) {
	http.Handle("/metrics", promhttp.Handler())
	port := fmt.Sprintf(":%d", configPort)
	log.Fatal(http.ListenAndServe(port, nil))
}
