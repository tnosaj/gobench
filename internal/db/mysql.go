package db

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	mysqlmigrate "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/sirupsen/logrus"
)

// ExecuteMySQL contains the connection and metrics to track executions
type ExecuteMySQL struct {
	Con     *sql.DB
	Metrics Metrics
}

// ExecStatement will execute a statement 's' and track it under the label 'l'
func (e ExecuteMySQL) ExecStatement(statement interface{}, label string) error {
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
func (e ExecuteMySQL) ExecInterfaceStatement(statement interface{}, label string) error {
	logrus.Debugf("will execut %q", statement)
	timer := prometheus.NewTimer(e.Metrics.DBRequestDuration.WithLabelValues(label))

	_, err := e.Con.Exec(stringInterfaceToSQLQuery(statement, label))
	if err != nil {
		e.Metrics.DBErrorRequests.WithLabelValues(label).Inc()
		return fmt.Errorf("could not execute %q with error %q", statement, err)
	}
	timer.ObserveDuration()
	return nil
}

// ExecStatementWithReturnBool will execute a statement 's' and return the resulting Boolean
func (e ExecuteMySQL) ExecStatementWithReturnBool(statement string) (bool, error) {
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
func (e ExecuteMySQL) ExecStatementWithReturnInt(statement string) (int, error) {
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
func (e ExecuteMySQL) ExecStatementWithReturnRow(statement, label string) (interface{}, error) {
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
func (e ExecuteMySQL) ExecDDLStatement(statement string) error {
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
func (e ExecuteMySQL) Ping() error {
	logrus.Debugf("will execut ping")

	if err := e.Con.Ping(); err != nil {
		logrus.Debugf("Failed to ping database: %s", err)
		return err
	}
	return nil
}

func connectMySQL(connectionInfo ConnectionInfo, poolsize int, metrics Metrics, tlsCerts TLSCerts) (*ExecuteMySQL, error) {
	logrus.Debugf("will connect to mysql")

	DSN := dsnFromConnectionInfo(connectionInfo)

	if tlsCerts.CaCertificate != "none" {
		logrus.Infof("using tls with cert: %s", tlsCerts.CaCertificate)
		rootCertPool := x509.NewCertPool()
		pem, err := os.ReadFile(tlsCerts.CaCertificate)
		if err != nil {
			logrus.Fatalf("failed to read certificate file: %s", err)
		}
		if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
			logrus.Fatalf("Failed to append PEM.")
		}

		clientCert := make([]tls.Certificate, 0, 1)
		var certs tls.Certificate

		if tlsCerts.ClientCertificate != "none" {
			certs, err = tls.LoadX509KeyPair(tlsCerts.ClientCertificate, tlsCerts.ClientKey)
			if err != nil {
				logrus.Fatal(err)
			}
		}

		clientCert = append(clientCert, certs)

		mysql.RegisterTLSConfig("custom", &tls.Config{
			RootCAs:            rootCertPool,
			Certificates:       clientCert,
			InsecureSkipVerify: true, // needed for self signed certs
		})
		DSN = DSN + "?tls=custom"

	} else {
		DSN = DSN + "?tls=skip-verify"
	}

	c, err := sql.Open("mysql", DSN)
	if err != nil {
		log.Fatalf("failed to open MySQL connection: %s", err)
	}

	c.SetMaxIdleConns(poolsize)
	c.SetMaxOpenConns(poolsize)
	c.SetConnMaxLifetime(360 * time.Second)
	return &ExecuteMySQL{Con: c, Metrics: metrics}, nil
}

func dsnFromConnectionInfo(connectionInfo ConnectionInfo) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		connectionInfo.User,
		connectionInfo.Password,
		connectionInfo.HostName,
		connectionInfo.Port,
		connectionInfo.DBName,
	)
}

func (e ExecuteMySQL) AutoMigrateUP(folder string) error {
	logrus.Debug("automatically migrating mysql up")
	driver, err := mysqlmigrate.WithInstance(e.Con, &mysqlmigrate.Config{
		MigrationsTable: "users_schema_migrations",
	})
	if err != nil {
		logrus.Errorf("Failed to create migration connection %s", err)
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s/mysql/", folder),
		"mysql",
		driver,
	)

	if err != nil {
		logrus.Errorf("Failed to initialize migration connection %s", err)
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logrus.Errorf("Failed to migrate db: %s", err)
		return err
	}
	logrus.Debug("succesfully migrated mysql up")
	return nil
}

func (e ExecuteMySQL) AutoMigrateDown(folder string) error {
	logrus.Debug("automatically migrating mysql down")
	driver, err := mysqlmigrate.WithInstance(e.Con, &mysqlmigrate.Config{
		MigrationsTable: "users_schema_migrations",
	})
	if err != nil {
		logrus.Errorf("Failed to create migration connection %s", err)
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s/mysql/", folder),
		"mysql",
		driver,
	)

	if err != nil {
		logrus.Errorf("Failed to initialize migration connection %s", err)
		return err
	}

	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		logrus.Errorf("Failed to migrate db: %s", err)
		return err
	}
	logrus.Debug("succesfully migrated mysql down")
	return nil
}

func (l *ExecuteMySQL) Shutdown(context context.Context) {
	logrus.Info("Shuttingdown longterm mysql server")
}
