package db

import (
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
)

// Row of work
type Row struct {
	ID  int
	K   int
	C   string
	Pad string
}

type TLSCerts struct {
	CaCertificate     string
	ClientCertificate string
	ClientKey         string
}

// Metrics contsins all metric types
type Metrics struct {
	DBRequestDuration *prometheus.HistogramVec
	DBErrorRequests   *prometheus.CounterVec
}

// ConnectionInfo contains all information needed to make a db connection
type ConnectionInfo struct {
	User     string
	Password string
	HostName string
	Port     string
	DBName   string
}

// DB Interface for all db operations
type DB interface {
	ExecStatement(statement, label string) error
	ExecStatementWithReturnBool(statement string) (bool, error)
	ExecStatementWithReturnInt(statement string) (int, error)
	ExecStatementWithReturnRow(statement, label string) (Row, error)
	ExecDDLStatement(statement string) error
	Ping() error

	// DB sepecific queries
	GetTableExists(dbName, tableName string) string
	Createable(dbName, tableName string) string
}

// Connect does the db magic connection
func Connect(db string, connectionInfo ConnectionInfo, poolsize int, tls TLSCerts) (DB, *sql.DB, error) {
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

	switch db {
	case "mysql":
		return connectMySQL(connectionInfo, poolsize, Metrics{
			DBRequestDuration: databaseRequestDuration,
			DBErrorRequests:   databaseErrorReuests,
		}, tls)
	case "postgres":
		return connectPostgreSQL(connectionInfo, poolsize, Metrics{
			DBRequestDuration: databaseRequestDuration,
			DBErrorRequests:   databaseErrorReuests,
		}, tls)
	default:
		return ExecuteNull{Metrics: Metrics{
			DBRequestDuration: databaseRequestDuration,
			DBErrorRequests:   databaseErrorReuests,
		}}, nil, nil

	}
}
