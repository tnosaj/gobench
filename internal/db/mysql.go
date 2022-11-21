package db

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

// ExecuteMySQL contains the connection and metrics to track executions
type ExecuteMySQL struct {
	Con     *sql.DB
	Metrics Metrics
}

// ExecStatement will execute a statement 's' and track it under the label 'l'
func (e ExecuteMySQL) ExecStatement(statement, label string) error {
	logrus.Debugf("will execut %q", statement)
	timer := prometheus.NewTimer(e.Metrics.DBRequestDuration.WithLabelValues(label))

	_, err := e.Con.Exec(statement)
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
func (e ExecuteMySQL) ExecStatementWithReturnRow(statement, label string) (Row, error) {
	logrus.Debugf("will execut %q", statement)

	timer := prometheus.NewTimer(e.Metrics.DBRequestDuration.WithLabelValues(label))
	var returnedRow Row

	q := e.Con.QueryRow(statement)

	if err := q.Scan(&returnedRow); err != nil {
		e.Metrics.DBErrorRequests.WithLabelValues(label).Inc()
		return Row{}, fmt.Errorf("query %q failed: %q", statement, err)
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

// GetTableExists returns a query to check if dbName.tableName exists
func (e ExecuteMySQL) GetTableExists(dbName, tableName string) string {
	return fmt.Sprintf("SELECT EXISTS(SELECT * FROM information_schema.tables "+
		"WHERE table_schema = '%s' AND table_name = '%s');", dbName, tableName)
}

// Createable returns a query to check if dbName.tableName exists
func (e ExecuteMySQL) Createable(dbName, tableName string) string {
	return fmt.Sprintf("CREATE TABLE %s.%s ( "+
		"id int(11) NOT NULL AUTO_INCREMENT, "+
		"k int(11) NOT NULL DEFAULT '0', "+
		"c char(120) NOT NULL DEFAULT '', "+
		"pad char(60) NOT NULL DEFAULT '', "+
		"PRIMARY KEY (id), "+
		"KEY k_1 (k)) ENGINE=InnoDB;", dbName, tableName)
}

func connectMySQL(connectionInfo ConnectionInfo, poolsize int, metrics Metrics, tlsCerts TLSCerts) (*ExecuteMySQL, *sql.DB, error) {
	logrus.Debugf("will connect to mysql")

	DSN := dsnFromConnectionInfo(connectionInfo)

	if tlsCerts.CaCertificate != "none" {
		logrus.Infof("using tls with cert: %s", tlsCerts.CaCertificate)
		rootCertPool := x509.NewCertPool()
		pem, err := ioutil.ReadFile(tlsCerts.CaCertificate)
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

	}

	c, err := sql.Open("mysql", DSN)
	if err != nil {
		log.Fatalf("failed to open MySQL connection: %s", err)
	}

	c.SetMaxIdleConns(poolsize)
	c.SetMaxOpenConns(poolsize)
	c.SetConnMaxLifetime(360 * time.Second)
	return &ExecuteMySQL{Con: c, Metrics: metrics}, c, nil
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
