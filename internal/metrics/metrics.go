package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"gitlab.otters.xyz/jason.tevnan/gobench/internal"
)

// RegisterPrometheusMetrics setsup metrics and returns them
func RegisterPrometheusMetrics() internal.Metrics {

	databaseRequestDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "database_request_duration_seconds",
		Help:    "Histogram for the runtime of a simple primary key get function.",
		Buckets: prometheus.LinearBuckets(0.01, 0.02, 75),
	}, []string{"query"})

	databaseErrorReuests := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "database_error_requests",
			Help: "The total number of failed requests",
		},
		[]string{"method"},
	)

	prometheus.MustRegister(databaseRequestDuration)
	prometheus.MustRegister(databaseErrorReuests)

	return internal.Metrics{DBRequestDuration: databaseRequestDuration, DBErrorRequests: databaseErrorReuests}
}
