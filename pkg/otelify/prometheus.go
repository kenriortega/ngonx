package otelify

import (
	"fmt"
	"net/http"

	"github.com/kenriortega/ngonx/pkg/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var MetricRequestLatencyProxy = promauto.NewHistogram(prometheus.HistogramOpts{
	Namespace: "ngonx",
	Name:      "request_latency_seconds",
	Help:      "Request Latency",
	Buckets:   prometheus.ExponentialBuckets(.0001, 2, 50),
})

func ExposeMetricServer(configPort int) {
	http.Handle("/metrics", promhttp.Handler())
	port := fmt.Sprintf(":%d", configPort)
	logger.LogError(http.ListenAndServe(port, nil).Error())
}
