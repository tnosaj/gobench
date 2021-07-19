package db

import (
	"gitlab.otters.xyz/jason.tevnan/gobench/internal"
)

// Connect does the db magic connection
func Connect(db string, connectionInfo internal.ConnectionInfo, poolsize int, metrics internal.Metrics, tls string) (internal.DB, error) {
	switch db {
	case "mysql":
		return connectMySQL(connectionInfo, poolsize, metrics, tls)
	case "postgres":
		return connectPostgreSQL(connectionInfo, poolsize, metrics, tls)
	default:
		return ExecuteNull{Metrics: metrics}, nil

	}
}
