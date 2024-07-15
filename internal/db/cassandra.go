package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/gocql/gocql"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/cassandra"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/sirupsen/logrus"
)

// ExecuteCassandra contains the connection and metrics to track executions
type ExecuteCassandra struct {
	Con     *gocql.Session
	Metrics Metrics
}

// ExecStatement will execute a statement 's' and track it under the label 'l'
func (e ExecuteCassandra) ExecStatement(statement interface{}, label string) error {
	logrus.Debugf("NOSQL DOES NOT SUPPORT THIS CALL: %q", statement)
	return nil
}

// ExecStatement will execute a statement 's' and track it under the label 'l'
func (e ExecuteCassandra) ExecInterfaceStatement(statement interface{}, label string) error {
	logrus.Tracef("will execut %q", statement)

	set := strings.Split(statement.(string), ",")
	timer := prometheus.NewTimer(e.Metrics.DBRequestDuration.WithLabelValues(label))
	var err error
	var (
		id  gocql.UUID
		k   gocql.UUID
		c   string
		pad string
	)
	switch label {
	case "read", "read-404":
		err = e.Con.Query(fmt.Sprintf("select id,k,c,pad from %s where id=?;", set[0]), set[1]).Consistency(gocql.One).Scan(&id, &k, &c, &pad)
	default:
		err = e.Con.Query(fmt.Sprintf("INSERT INTO %s (id, k, c, pad) VALUES (?,?,?,?);", set[0]), set[1], set[1], set[2], set[2]).Exec()
	}
	timer.ObserveDuration()
	if err != nil {
		if errors.Is(err, gocql.ErrNotFound) {
			logrus.Tracef("error: %s", err.Error())
			// not an error, but we want to count it
			e.Metrics.DBErrorRequests.WithLabelValues(fmt.Sprintf("%s-404", label)).Inc()
			return nil
		}
		logrus.Warnf("error: %s", err.Error())
		e.Metrics.DBErrorRequests.WithLabelValues(label).Inc()
		return fmt.Errorf("could not execute %q with error %q", statement, err)
	}
	return nil
}

// ExecStatementWithReturnBool will execute a statement 's' and return the resulting Boolean
func (e ExecuteCassandra) ExecStatementWithReturnBool(statement string) (bool, error) {
	logrus.Debugf("NOSQL DOES NOT SUPPORT THIS CALL: %q", statement)

	var returnedBoolean bool
	return returnedBoolean, nil
}

// ExecStatementWithReturnInt will execute a statement 's' and return the resulting Int
func (e ExecuteCassandra) ExecStatementWithReturnInt(statement string) (int, error) {
	logrus.Debugf("NOSQL DOES NOT SUPPORT THIS CALL: %q", statement)
	var returnedInt int
	return returnedInt, nil
}

// ExecStatementWithReturnRow will execute a statement 's', track it under the label 'l' and return the resulting Row
func (e ExecuteCassandra) ExecStatementWithReturnRow(statement, label string) (interface{}, error) {
	logrus.Debugf("NOSQL DOES NOT SUPPORT THIS CALL: %q", statement)
	var returnedRow interface{}
	return returnedRow, nil
}

// ExecDDLStatement will execute a statement 's' as a DDL
func (e ExecuteCassandra) ExecDDLStatement(statement string) error {
	logrus.Debugf("NOSQL DOES NOT SUPPORT THIS CALL: %q", statement)
	return nil
}

// Ping checks if the db is up
func (e ExecuteCassandra) Ping() error {
	logrus.Debug("NOSQL DOES NOT SUPPORT ping CALL")
	return nil
}

func connectCassandra(connectionInfo ConnectionInfo, poolsize int, metrics Metrics, tlsCerts TLSCerts) (*ExecuteCassandra, error) {
	logrus.Debugf("will connect to Cassandra")

	hostsList := strings.Split(connectionInfo.HostName, ",")

	cluster := gocql.NewCluster(hostsList...)
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: connectionInfo.User,
		Password: connectionInfo.Password,
	}
	// hardcoded
	//cluster.Timeout = 10
	//cluster.ConnectTimeout = 10
	//cluster.WriteTimeout = 10

	cluster.Consistency = gocql.Quorum
	cluster.Keyspace = connectionInfo.DBName
	cluster.NumConns = poolsize

	if tlsCerts.CaCertificate != "none" {
		cluster.SslOpts = &gocql.SslOptions{
			EnableHostVerification: true,
		}

	}

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalln("Failed to connect to the server cluster: ", err)
	}

	return &ExecuteCassandra{Con: session, Metrics: metrics}, nil
}

func (e ExecuteCassandra) AutoMigrateUP(folder string) error {
	logrus.Debug("automatically migrating Cassandra up")
	driver, err := cassandra.WithInstance(e.Con, &cassandra.Config{
		MigrationsTable: "schema_migrations",
		KeyspaceName:    "sbtest",
	})
	if err != nil {
		logrus.Errorf("Failed to create migration connection %s", err)
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s/cassandra/", folder),
		"cassandra",
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
	logrus.Debug("succesfully migrated Cassandra up")
	return nil
}

func (e ExecuteCassandra) AutoMigrateDown(folder string) error {
	logrus.Debug("automatically migrating Cassandra down")
	driver, err := cassandra.WithInstance(e.Con, &cassandra.Config{
		MigrationsTable: "schema_migrations",
		KeyspaceName:    "sbtest",
	})
	if err != nil {
		logrus.Errorf("Failed to create migration connection %s", err)
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s/cassandra/", folder),
		"cassandra",
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
	logrus.Debug("succesfully migrated Cassandra down")
	return nil
}

func (e *ExecuteCassandra) Shutdown(context context.Context) {
	logrus.Info("Shuttingdown Cassandra connections")
	e.Con.Close()
}
