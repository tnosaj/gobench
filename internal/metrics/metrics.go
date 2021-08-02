package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"gitlab.otters.xyz/jason.tevnan/gobench/internal"
)

// RegisterPrometheusMetrics setsup metrics and returns them
func RegisterPrometheusMetrics() internal.Metrics {

	mysqlRequestDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "mysl_request_duration_seconds",
		Help:    "Histogram for the runtime of a simple primary key get function.",
		Buckets: prometheus.LinearBuckets(0.00, 0.002, 50),
	}, []string{"query"})

	mysqlErrorReuests := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "mysql_error_requests",
			Help: "The total number of failed requests",
		},
		[]string{"method"},
	)

	prometheus.MustRegister(mysqlRequestDuration)
	prometheus.MustRegister(mysqlErrorReuests)

	return internal.Metrics{DBRequestDuration: mysqlRequestDuration, DBErrorRequests: mysqlErrorReuests}
}
