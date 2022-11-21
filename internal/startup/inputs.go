package startup

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
	"gitlab.otters.xyz/jason.tevnan/gobench/internal"
	"gitlab.otters.xyz/jason.tevnan/gobench/internal/helper"
	"gitlab.otters.xyz/jason.tevnan/gobench/pkg/args"
)

// EvaluateInputs checks for sanity
func EvaluateInputs() (internal.Settings, error) {
	var s internal.Settings

	flag.BoolVar(&s.Debug, "v", false, "Enable verbose debugging output")

	flag.StringVar(&s.Action, "a", "unset", "Perform this action, must be one of the following: prepare,run,cleanup")
	flag.StringVar(&s.TableName, "n", "sbtest", "Base table name")
	flag.StringVar(&s.Port, "p", "8080", "Starts server on this port")
	flag.StringVar(&s.Strategy, "strategy", "simple", "Strategy to use")
	flag.StringVar(&s.TLSCerts.CaCertificate, "cacert", "none", "/path/to/ca-certificate.pem if using tls")
	flag.StringVar(&s.TLSCerts.ClientCertificate, "clientcert", "none", "/path/to/client-certificate.pem if using tls")
	flag.StringVar(&s.TLSCerts.ClientKey, "clientkey", "none", "/path/to/client-key.key if using tls")
	flag.StringVar(&s.DurationType, "duration", "events", "Duratation type - events|seconds")

	flag.IntVar(&s.Concurrency, "t", 1, "Concurrent number of threads to run")
	flag.IntVar(&s.Connectionpoolsize, "c", 20, "Connection pool size")
	flag.IntVar(&s.Initialdatasize, "i", 1000, "Initial size to populate")
	flag.IntVar(&s.Duration, "d", 1000, "The number of events to process")
	flag.IntVar(&s.Rate, "r", 0, "requests per second - 0 to disable rate limiting")

	var split string
	flag.StringVar(&split, "s", "90:10", "Read:Write split, seperated by a colon")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s: [flags] command [command args…]\n", os.Args[0])
		flag.PrintDefaults()
		printRequiredEnvVars()
	}

	flag.Parse()

	setupLogger(s.Debug)

	rws, err := args.ParseReadWriteSplit(split)
	if err != nil {
		return internal.Settings{}, err
	}
	s.ReadWriteSplit = rws

	s.DBConnectionInfo.User, err = getEnvVar("DBUSER")
	if err != nil {
		return internal.Settings{}, err
	}
	s.DBConnectionInfo.Password, err = getEnvVar("DBPASSWORD")
	if err != nil {
		return internal.Settings{}, err
	}
	s.DBConnectionInfo.HostName, err = getEnvVar("DBHOSTNAME")
	if err != nil {
		return internal.Settings{}, err
	}
	s.DBConnectionInfo.DBName, err = getEnvVar("DBNAME")
	if err != nil {
		return internal.Settings{}, err
	}
	s.DBConnectionInfo.Port, err = getEnvVar("DBPORT")
	if err != nil {
		return internal.Settings{}, err
	}
	s.DB, err = getEnvVar("DBENGINE")
	if err != nil {
		return internal.Settings{}, err
	}

	acceptedDBs := map[string]bool{
		"mysql":    true,
		"postgres": true,
		"nulldb":   true,
	}
	if !acceptedDBs[s.DB] {
		return internal.Settings{}, fmt.Errorf("Unknown db engine specified: %q", s.DB)
	}

	for _, tlsfile := range []string{s.TLSCerts.CaCertificate, s.TLSCerts.ClientCertificate, s.TLSCerts.ClientKey} {
		if tlsfile != "none" {
			_, err := ioutil.ReadFile(tlsfile)
			if err != nil {
				return internal.Settings{}, fmt.Errorf("%q set, but unreadable path: %q", tlsfile, err)
			}
		}
	}

	acceptedActions := map[string]bool{
		"prepare": true,
		"run":     true,
		"cleanup": true,
	}
	if !acceptedActions[s.Action] {
		return internal.Settings{}, fmt.Errorf("Action %q is not one of the following: prepare, run, cleanup", s.Action)
	}

	acceptedDurationTypes := map[string]bool{
		"events":  true,
		"seconds": true,
	}
	if !acceptedDurationTypes[s.DurationType] {
		return internal.Settings{}, fmt.Errorf("DurationType %q is not one of the following: events, seconds", s.DurationType)
	}

	s.Randomizer = helper.NewRandomizer()

	return s, nil
}

func getEnvVar(envVarName string) (string, error) {
	check := os.Getenv(envVarName)
	if check == "" {
		printRequiredEnvVars()
		return "", fmt.Errorf("Missing env var %q", envVarName)
	}
	return check, nil
}

func setupLogger(debug bool) {
	//logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
	logrus.Debug("Configured logger")
}

func printRequiredEnvVars() {
	fmt.Println("\nRequired Environment variables:")
	fmt.Println("  DBUSER - user to connect as")
	fmt.Println("  DBPASSWORD - password to connect with")
	fmt.Println("  DBHOSTNAME - instance to connect to")
	fmt.Println("  DBNAME - db name to connect to on the instance")
	fmt.Println("  DBPORT - port of the db instance")
	fmt.Println("  DBENGINE - type of db")
}
