package db

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"
	"strconv"
	"strings"

	as "github.com/aerospike/aerospike-client-go/v7"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/sirupsen/logrus"
)

// ExecuteAerospike contains the connection and metrics to track executions
type ExecuteAerospike struct {
	Con     *as.Client
	Metrics Metrics
}

type AerospikeData struct {
	namespace string
	setname   string
	key       string
	binmap    string
}

// ExecStatement will execute a statement 's' and track it under the label 'l'
func (e ExecuteAerospike) ExecStatement(statement, label string) error {
	logrus.Debugf("will execut %q", statement)
	data := stringToAreospikeData(statement)
	key, _ := as.NewKey(data.namespace, data.setname, data.key)
	timer := prometheus.NewTimer(e.Metrics.DBRequestDuration.WithLabelValues(label))
	switch label {
	case "read":
		e.Con.Get(nil, key)
	default:
		bin := as.BinMap{"data": data.binmap}
		e.Con.Put(nil, key, bin)
	}
	// _, err := e.Con.Exec(statement)
	// if err != nil {
	// 	e.Metrics.DBErrorRequests.WithLabelValues(label).Inc()
	// 	return fmt.Errorf("could not execute %q with error %q", statement, err)
	// }
	timer.ObserveDuration()
	return nil
}

// ExecStatementWithReturnBool will execute a statement 's' and return the resulting Boolean
func (e ExecuteAerospike) ExecStatementWithReturnBool(statement string) (bool, error) {
	logrus.Debugf("will execut %q", statement)
	var returnedBoolean bool

	// q := e.Con.QueryRow(statement)

	// if err := q.Scan(&returnedBoolean); err != nil {

	// 	return false, fmt.Errorf("query %q failed: %q", statement, err)
	// }
	logrus.Debugf("returning %t", returnedBoolean)
	return returnedBoolean, nil
}

// ExecStatementWithReturnInt will execute a statement 's' and return the resulting Int
func (e ExecuteAerospike) ExecStatementWithReturnInt(statement string) (int, error) {
	logrus.Debugf("will execut %q", statement)
	var returnedInt int

	// q := e.Con.QueryRow(statement)

	// if err := q.Scan(&returnedInt); err != nil {
	// 	return 0, fmt.Errorf("query %q failed: %q", statement, err)
	// }

	logrus.Debugf("returning %d", returnedInt)
	return returnedInt, nil
}

// ExecStatementWithReturnRow will execute a statement 's', track it under the label 'l' and return the resulting Row
func (e ExecuteAerospike) ExecStatementWithReturnRow(statement, label string) (interface{}, error) {
	logrus.Debugf("will execut %q", statement)

	timer := prometheus.NewTimer(e.Metrics.DBRequestDuration.WithLabelValues(label))
	var returnedRow interface{}

	// q := e.Con.QueryRow(statement)

	// if err := q.Scan(&returnedRow); err != nil {
	// 	e.Metrics.DBErrorRequests.WithLabelValues(label).Inc()
	// 	return returnedRow, fmt.Errorf("query %q failed: %q", statement, err)
	// }
	timer.ObserveDuration()

	logrus.Debugf("returning %+v", returnedRow)

	return returnedRow, nil
}

// ExecDDLStatement will execute a statement 's' as a DDL
func (e ExecuteAerospike) ExecDDLStatement(statement string) error {
	logrus.Debugf("will execut %q", statement)
	// not implemented
	return nil
}

// Ping checks if the db is up
func (e ExecuteAerospike) Ping() error {
	logrus.Debugf("will execut ping")
	// not implemented
	return nil
}

func connectAerospike(connectionInfo ConnectionInfo, metrics Metrics, tlsCerts TLSCerts) (*ExecuteAerospike, error) {
	logrus.Debugf("will connect to aerospike")

	var client *as.Client
	port, err := strconv.Atoi(connectionInfo.Port)

	if err != nil {
		logrus.Fatal(err)
	}
	clientPolicy := as.NewClientPolicy()

	// hard coded
	clientPolicy.IdleTimeout = 10
	clientPolicy.ConnectionQueueSize = 100

	multiPolicy := as.NewMultiPolicy()
	multiPolicy.MaxRecords = 1

	basePolicy := as.NewPolicy()
	basePolicy.SendKey = true

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

		clientPolicy.TlsConfig = &tls.Config{
			RootCAs:            rootCertPool,
			Certificates:       clientCert,
			InsecureSkipVerify: true, // needed for self signed certs
		}

	}

	client, err = as.NewClientWithPolicy(clientPolicy, connectionInfo.HostName, port)
	client.DefaultPolicy = basePolicy
	client.DefaultQueryPolicy.MultiPolicy = *multiPolicy
	client.DefaultScanPolicy.MultiPolicy = *multiPolicy

	if err != nil {
		log.Fatalln("Failed to connect to the server cluster: ", err)
	}

	return &ExecuteAerospike{Con: client, Metrics: metrics}, nil
}

func (e ExecuteAerospike) AutoMigrateUP(folder string) error {
	logrus.Debug("automatically migrating aerospike up")
	// not implemented
	logrus.Debug("succesfully migrated aerospike up")
	return nil
}

func (e ExecuteAerospike) AutoMigrateDown(folder string) error {
	logrus.Debug("automatically migrating aerospike down")
	// not implemented
	logrus.Debug("succesfully migrated aerospike down")
	return nil
}

func (e *ExecuteAerospike) Shutdown(context context.Context) {
	logrus.Info("Shuttingdown aerospike connections")
}

func stringToAreospikeData(s string) AerospikeData {
	var a AerospikeData
	set := strings.Split(s, ",")
	a.namespace = set[0]
	a.setname = set[1]
	a.key = set[2]
	if len(set) > 3 {
		a.binmap = set[3]
	}
	return a
}
