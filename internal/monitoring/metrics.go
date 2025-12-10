package monitoring

import (
	"log"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

var (
	customRegistry = prometheus.NewRegistry()

	RequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_requests_total",
			Help: "Total number of gRPC requests",
		},
		[]string{"service", "method", "code"},
	)

	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "grpc_request_duration_milliseconds",
			Help:    "gRPC request duration in milliseconds",
			Buckets: []float64{1, 5, 10, 25, 50, 100, 250, 500, 1000, 5000},
		},
		[]string{"service", "method"},
	)

	RequestsInFlight = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "grpc_requests_in_flight",
			Help: "Current number of gRPC requests in flight",
		},
		[]string{"service", "method"},
	)

	initOnce sync.Once
)

func InitMetrics() {
	initOnce.Do(func() {
		customRegistry.MustRegister(RequestsTotal)
		customRegistry.MustRegister(RequestDuration)
		customRegistry.MustRegister(RequestsInFlight)

		customRegistry.MustRegister(collectors.NewGoCollector())
		customRegistry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

		log.Println("Metrics initialized with custom registry")
	})
}
