package metric

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
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
)

func init() {
	prometheus.MustRegister(CountersByRoute)
}

func ExposeMetricServer(configPort int) {
	http.Handle("/metrics", promhttp.Handler())
	port := fmt.Sprintf(":%d", configPort)
	log.Fatal(http.ListenAndServe(port, nil))
}
