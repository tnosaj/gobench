package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/sirupsen/logrus"
)

// ExecutePostSQL contains the connection and metrics to track executions
type ExecutePostSQL struct {
	Con     *sql.DB
	Metrics Metrics
}

// ExecStatement will execute a statement 's' and track it under the label 'l'
func (e ExecutePostSQL) ExecStatement(statement interface{}, label string) error {
	logrus.Debugf("will execut %q", statement)
	timer := prometheus.NewTimer(e.Metrics.DBRequestDuration.WithLabelValues(label))

	_, err := e.Con.Exec(statement.(string))
	if err != nil {
		e.Metrics.DBErrorRequests.WithLabelValues(label).Inc()
		return fmt.Errorf("could not execute %q with error %q", statement, err)
	}
	timer.ObserveDuration()
	return nil
}

// ExecStatement will execute a statement 's' and track it under the label 'l'
func (e ExecutePostSQL) ExecInterfaceStatement(statement interface{}, label string) error {
	logrus.Tracef("will execut %q", statement)
	timer := prometheus.NewTimer(e.Metrics.DBRequestDuration.WithLabelValues(label))

	_, err := e.Con.Exec(stringInterfaceToPostgreSQLQuery(statement, label))
	if err != nil {
		e.Metrics.DBErrorRequests.WithLabelValues(label).Inc()
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
func (e ExecutePostSQL) ExecStatementWithReturnRow(statement, label string) (interface{}, error) {
	logrus.Debugf("will execut %q", statement)

	timer := prometheus.NewTimer(e.Metrics.DBRequestDuration.WithLabelValues(label))
	var returnedRow interface{}

	q := e.Con.QueryRow(statement)

	if err := q.Scan(&returnedRow); err != nil {
		e.Metrics.DBErrorRequests.WithLabelValues(label).Inc()
		return returnedRow, fmt.Errorf("query %q failed: %q", statement, err)
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

func connectPostgreSQL(connectionInfo ConnectionInfo, poolsize int, metrics Metrics, tlsCerts TLSCerts) (*ExecutePostSQL, error) {
	logrus.Debugf("will connect to postgres")

	var psqlInfo string
	if tlsCerts.CaCertificate != "none" {
		logrus.Infof("using tls with cert: %s", tlsCerts.CaCertificate)

		_, err := os.ReadFile(tlsCerts.CaCertificate)
		if err != nil {
			logrus.Fatalf("failed to read certificate file: %s", err)
		}

		psqlInfo = psqlInfo + "sslmode=require sslrootcert=" + tlsCerts.CaCertificate
		if tlsCerts.ClientCertificate != "none" {
			psqlInfo = psqlInfo + " sslkey=" + tlsCerts.ClientKey + " sslcert=" + tlsCerts.ClientCertificate
		}

	} else {
		psqlInfo = psqlInfo + "sslmode=disable"
	}
	//
	// moved this to the end because I was
	// having issues with passwords with
	// special chars messing with ssl settings
	//
	psqlInfo = psqlInfo + " " + psqlInfoFromConnectionInfo(connectionInfo)
	c, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("failed to open PostgreSQL connection: %s", err)
	}

	c.SetMaxIdleConns(poolsize)
	c.SetMaxOpenConns(poolsize)
	c.SetConnMaxLifetime(360 * time.Second)
	return &ExecutePostSQL{Con: c, Metrics: metrics}, nil
}

func psqlInfoFromConnectionInfo(connectionInfo ConnectionInfo) string {
	return fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s",
		connectionInfo.User,
		connectionInfo.Password,
		connectionInfo.HostName,
		connectionInfo.Port,
		connectionInfo.DBName,
	)
}

func (e ExecutePostSQL) AutoMigrateUP(folder string) error {
	logrus.Debug("automatically migrating postgres")
	driver, err := postgres.WithInstance(e.Con, &postgres.Config{
		MigrationsTable: "schema_migrations",
	})
	if err != nil {
		logrus.Errorf("Failed to create migration connection %s", err)
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s/postgres", folder),
		"postgres", driver,
	)
	if err != nil {
		logrus.Errorf("Failed to initialize migration connection %s", err)
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logrus.Errorf("Failed to migrate db: %s", err)
		return err
	}
	logrus.Debug("successfully migrated postgres")

	return nil
}
func (e ExecutePostSQL) AutoMigrateDown(folder string) error {
	logrus.Debug("automatically migrating postgres")
	driver, err := postgres.WithInstance(e.Con, &postgres.Config{
		MigrationsTable: "schema_migrations",
	})
	if err != nil {
		logrus.Errorf("Failed to create migration connection %s", err)
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s/postgres", folder),
		"postgres", driver,
	)
	if err != nil {
		logrus.Errorf("Failed to initialize migration connection %s", err)
		return err
	}

	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		logrus.Errorf("Failed to migrate db: %s", err)
		return err
	}
	logrus.Debug("successfully migrated postgres")

	return nil
}

func (l *ExecutePostSQL) Shutdown(context context.Context) {
	logrus.Info("Shuttingdown longterm postgres server")
}

func stringInterfaceToPostgreSQLQuery(s interface{}, label string) string {
	set := strings.Split(s.(string), ",")

	switch label {
	case "read":
		return fmt.Sprintf("select id,k,c,pad from %s where id='%s';", set[0], set[1])
	}
	return fmt.Sprintf("INSERT INTO %s(id, k, c , pad) VALUES ('%s','%s','%s','%s');", set[0], set[1], set[2], set[2], set[2])

}
