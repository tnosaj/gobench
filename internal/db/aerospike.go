package db

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
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
func (e ExecuteAerospike) ExecStatement(statement interface{}, label string) error {
	logrus.Debugf("NOSQL DOES NOT SUPPORT THIS CALL: %q", statement)
	return nil
}

// ExecStatement will execute a statement 's' and track it under the label 'l'
func (e ExecuteAerospike) ExecInterfaceStatement(statement interface{}, label string) error {
	logrus.Tracef("will execut %q", statement)
	data := stringToAerospikeData(statement)
	key, _ := as.NewKey(data.namespace, data.setname, data.key)
	timer := prometheus.NewTimer(e.Metrics.DBRequestDuration.WithLabelValues(label))
	var err as.Error
	switch label {
	case "read", "read-404":
		_, err = e.Con.Get(nil, key)
	default:
		bin := as.BinMap{"data": data.binmap}
		err = e.Con.Put(nil, key, bin)
	}
	if err != nil {
		e.Metrics.DBErrorRequests.WithLabelValues(label).Inc()
		return fmt.Errorf("could not execute %q with error %q", statement, err)
	}
	timer.ObserveDuration()
	return nil
}

// ExecStatementWithReturnBool will execute a statement 's' and return the resulting Boolean
func (e ExecuteAerospike) ExecStatementWithReturnBool(statement string) (bool, error) {
	logrus.Debugf("NOSQL DOES NOT SUPPORT THIS CALL: %q", statement)

	var returnedBoolean bool
	return returnedBoolean, nil
}

// ExecStatementWithReturnInt will execute a statement 's' and return the resulting Int
func (e ExecuteAerospike) ExecStatementWithReturnInt(statement string) (int, error) {
	logrus.Debugf("NOSQL DOES NOT SUPPORT THIS CALL: %q", statement)
	var returnedInt int
	return returnedInt, nil
}

// ExecStatementWithReturnRow will execute a statement 's', track it under the label 'l' and return the resulting Row
func (e ExecuteAerospike) ExecStatementWithReturnRow(statement, label string) (interface{}, error) {
	logrus.Debugf("NOSQL DOES NOT SUPPORT THIS CALL: %q", statement)
	var returnedRow interface{}
	return returnedRow, nil
}

// ExecDDLStatement will execute a statement 's' as a DDL
func (e ExecuteAerospike) ExecDDLStatement(statement string) error {
	logrus.Debugf("NOSQL DOES NOT SUPPORT THIS CALL: %q", statement)
	return nil
}

// Ping checks if the db is up
func (e ExecuteAerospike) Ping() error {
	logrus.Debug("NOSQL DOES NOT SUPPORT THIS CALL")
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
	if err != nil {
		log.Fatalln("Failed to connect to the server cluster: ", err)
	}
	client.DefaultPolicy = basePolicy
	client.DefaultQueryPolicy.MultiPolicy = *multiPolicy
	client.DefaultScanPolicy.MultiPolicy = *multiPolicy

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
	e.Con.Close()
}

func stringToAerospikeData(s interface{}) AerospikeData {
	var a AerospikeData
	set := strings.Split(s.(string), ",")
	a.namespace = set[0]
	a.setname = set[0]
	a.key = set[1]
	if len(set) > 2 {
		a.binmap = set[2]
	}
	return a
}
