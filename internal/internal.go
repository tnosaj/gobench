package internal

import (
	"github.com/prometheus/client_golang/prometheus"
	"gitlab.otters.xyz/jason.tevnan/gobench/pkg/args"
)

// ConnectionInfo contains all information needed to make a db connection
type ConnectionInfo struct {
	User     string
	Password string
	HostName string
	Port     string
	DBName   string
}

// Settings contains global settings
type Settings struct {
	Concurrency        int
	Connectionpoolsize int
	Initialdatasize    int
	Duration           int
	Rate               int

	Debug bool

	Action           string
	DB               string
	DBConnectionInfo ConnectionInfo
	DurationType     string
	Port             string
	Strategy         string
	TableName        string

	TLSCerts TLSCerts

	DBInterface    DB
	ReadWriteSplit args.ReadWriteSplit
	Randomizer     Random
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

// select * from sbtest1 where id=1\G
// *************************** 1. row ***************************
// id: 1
// k: 119929
// c: 97910776665-24067652665-16947278706-50317850100-81651631236-65011289994-96675195764-77928027135-17068252874-45254863248
// pad: 73947976563-21217985552-23922161426-64060646662-89232940667
//
// id int(11) NOT NULL AUTO_INCREMENT
// k int(11) NOT NULL DEFAULT '0',
// c char(120) NOT NULL DEFAULT ''
// pad char(60) NOT NULL DEFAULT ''

// Row of work
type Row struct {
	ID  int
	K   int
	C   string
	Pad string
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

// ExecutionStrategy defines what queries are run how
type ExecutionStrategy interface {
	CreateCommand() (string, string)
}

// ExecutionType defines how long queries are run
type ExecutionType interface {
	Run()
}

// Random interface to help with testing
type Random interface {
	Intn(int) int
}
