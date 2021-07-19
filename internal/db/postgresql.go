package db

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	_ "github.com/lib/pq"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"gitlab.otters.xyz/jason.tevnan/gobench/internal"
)

// ExecutePostSQL contains the connection and metrics to track executions
type ExecutePostSQL struct {
	Con     *sql.DB
	Metrics internal.Metrics
}

// ExecStatement will execute a statement 's' and track it under the label 'l'
func (e ExecutePostSQL) ExecStatement(statement, label string) error {
	logrus.Debugf("will execut %q", statement)
	timer := prometheus.NewTimer(e.Metrics.DBRequestDuration.WithLabelValues(label))

	_, err := e.Con.Exec(statement)
	if err != nil {
		e.Metrics.DBErrorReuests.WithLabelValues(label).Inc()
		return fmt.Errorf("could not execute %q with error %q", statement, err)
	}
	timer.ObserveDuration()
	return nil
}

// ExecStatementWithReturnBool will execute a statement 's' and return the resulting Boolean
func (e ExecutePostSQL) ExecStatementWithReturnBool(statement string) (bool, error) {
	logrus.Debugf("will execut %q", statement)
	var returnedBoolean bool

	q := e.Con.QueryRow(statement)

	if err := q.Scan(&returnedBoolean); err != nil {

		return false, fmt.Errorf("query %q failed: %q", statement, err)
	}
	logrus.Debugf("returning %t", returnedBoolean)
	return returnedBoolean, nil
}

// ExecStatementWithReturnInt will execute a statement 's' and return the resulting Int
func (e ExecutePostSQL) ExecStatementWithReturnInt(statement string) (int, error) {
	logrus.Debugf("will execut %q", statement)
	var returnedInt int

	q := e.Con.QueryRow(statement)

	if err := q.Scan(&returnedInt); err != nil {
		return 0, fmt.Errorf("query %q failed: %q", statement, err)
	}

	logrus.Debugf("returning %d", returnedInt)
	return returnedInt, nil
}

// ExecStatementWithReturnRow will execute a statement 's', track it under the label 'l' and return the resulting Row
func (e ExecutePostSQL) ExecStatementWithReturnRow(statement, label string) (internal.Row, error) {
	logrus.Debugf("will execut %q", statement)

	timer := prometheus.NewTimer(e.Metrics.DBRequestDuration.WithLabelValues(label))
	var returnedRow internal.Row

	q := e.Con.QueryRow(statement)

	if err := q.Scan(&returnedRow); err != nil {
		e.Metrics.DBErrorReuests.WithLabelValues(label).Inc()
		return internal.Row{}, fmt.Errorf("query %q failed: %q", statement, err)
	}
	timer.ObserveDuration()

	logrus.Debugf("returning %+v", returnedRow)

	return returnedRow, nil
}

// ExecDDLStatement will execute a statement 's' as a DDL
func (e ExecutePostSQL) ExecDDLStatement(statement string) error {
	logrus.Debugf("will execut %q", statement)

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	_, err := e.Con.ExecContext(ctx, statement)
	if err != nil {
		return fmt.Errorf("could not run %q with error %q", statement, err)
	}
	return nil
}

// Ping checks if the db is up
func (e ExecutePostSQL) Ping() error {
	logrus.Debugf("will execut ping")

	if err := e.Con.Ping(); err != nil {
		logrus.Debugf("Failed to ping database: %s", err)
		return err
	}
	return nil
}

// GetTableExists returns a query to check if dbName.tableName exists
func (e ExecutePostSQL) GetTableExists(dbName, tableName string) string {
	return fmt.Sprintf("SELECT EXISTS ( "+
		"SELECT FROM information_schema.tables "+
		"WHERE  table_catalog = '%s' "+
		"AND    table_name   = '%s' );", dbName, tableName)
}

// Createable returns a query to check if dbName.tableName exists
func (e ExecutePostSQL) Createable(dbName, tableName string) string {
	return fmt.Sprintf("CREATE TABLE %s ( "+
		" id SERIAL PRIMARY KEY, "+
		" k integer NOT NULL DEFAULT '0', "+
		" c VARCHAR NOT NULL DEFAULT '', "+
		" pad VARCHAR NOT NULL DEFAULT ''); "+
		"CREATE INDEX k_idx ON %s (k);", tableName, tableName)
}

func connectPostgreSQL(connectionInfo internal.ConnectionInfo, poolsize int, metrics internal.Metrics, tlsCert string) (*ExecutePostSQL, error) {
	logrus.Debugf("will connect to postgres")

	psqlInfo := psqlInfoFromConnectionInfo(connectionInfo)
	if tlsCert != "none" {
		logrus.Infof("using tls with cert: %s", tlsCert)

		_, err := ioutil.ReadFile(tlsCert)
		if err != nil {
			log.Fatalf("failed to read certificate file: %s", err)
		}

		psqlInfo = psqlInfo + " sslmode=verify-full slrootcert=" + tlsCert

	} else {
		psqlInfo = psqlInfo + " sslmode=disable"
	}

	c, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("failed to open PostgreSQL connection: %s", err)
	}

	c.SetMaxIdleConns(poolsize)
	c.SetMaxOpenConns(poolsize)
	c.SetConnMaxLifetime(360 * time.Second)
	return &ExecutePostSQL{Con: c, Metrics: metrics}, nil
}

func psqlInfoFromConnectionInfo(connectionInfo internal.ConnectionInfo) string {
	return fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s",
		connectionInfo.User,
		connectionInfo.Password,
		connectionInfo.HostName,
		connectionInfo.Port,
		connectionInfo.DBName,
	)
}
