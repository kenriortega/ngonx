package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/kenriortega/ngonx/pkg/logger"
	"github.com/kenriortega/ngonx/pkg/otelify"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// global vars...gasp!
var addr = "127.0.0.1:8000"
var tracer trace.Tracer
var httpClient http.Client

func main() {
	flush := otelify.InitProvider(
		"example",
		"v0.4.5",
		"test",
		"0.0.0.0:55680",
	)
	defer flush()

	// initiate globals
	tracer = otel.Tracer("example-app")
	httpClient = http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	// create and start server
	server := instrumentedServer(handler)

	fmt.Println("listening...")
	log.Fatal(server.ListenAndServe())
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	longRunningProcess(ctx)

	// check cache
	if shouldExecute(40) {
		url := "http://" + addr + "/"

		resp, err := instrumentedGet(ctx, url)
		defer resp.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	// query database
	if shouldExecute(40) {
		url := "http://" + addr + "/"

		resp, err := instrumentedGet(ctx, url)
		defer resp.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func shouldExecute(percent int) bool {
	return rand.Int()%100 < percent
}

func longRunningProcess(ctx context.Context) {
	ctx, sp := tracer.Start(ctx, "Long Running Process")
	defer sp.End()

	time.Sleep(time.Millisecond * 50)
	sp.AddEvent("halfway done!")
	time.Sleep(time.Millisecond * 50)
}

/***
Server
***/
func instrumentedServer(handler http.HandlerFunc) *http.Server {
	// OpenMetrics handler : metrics and exemplars
	omHandleFunc := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		handler.ServeHTTP(w, r)

		ctx := r.Context()
		traceID := trace.SpanContextFromContext(ctx).TraceID().String()

		otelify.MetricRequestLatencyProxy.(prometheus.ExemplarObserver).ObserveWithExemplar(
			time.Since(start).Seconds(), prometheus.Labels{"traceID": traceID},
		)

		// log the trace id with other fields so we can discover traces through logs
		logger.LogInfo(
			"http request",
			zap.String("traceID", traceID),
			zap.String("path", r.URL.Path),
			zap.Duration("latency", time.Since(start)),
		)
	}

	// OTel handler : traces
	otelHandler := otelhttp.NewHandler(http.HandlerFunc(omHandleFunc), "http")

	r := mux.NewRouter()
	r.Handle("/", otelHandler)
	r.Handle("/metrics", promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{
		EnableOpenMetrics: true,
	}))

	return &http.Server{
		Handler: r,
		Addr:    "0.0.0.0:8000",
	}
}

/***
Client
***/
func instrumentedGet(ctx context.Context, url string) (*http.Response, error) {
	// create http request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		panic(err)
	}

	return httpClient.Do(req)
}

func handleErr(err error, message string) {
	if err != nil {
		panic(fmt.Sprintf("%s: %s", err, message))
	}
}
