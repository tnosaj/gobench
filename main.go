package main

import (
	"math/rand"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"gitlab.otters.xyz/jason.tevnan/gobench/internal/db"
	"gitlab.otters.xyz/jason.tevnan/gobench/internal/metrics"
	"gitlab.otters.xyz/jason.tevnan/gobench/internal/startup"
	"gitlab.otters.xyz/jason.tevnan/gobench/internal/work"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	s, err := startup.EvaluateInputs()
	if err != nil {
		log.Fatalf("could not evaluate inputs: %q", err)
	}

	connectionInterface, err := db.Connect(s.DB,
		s.DBConnectionInfo,
		s.Connectionpoolsize,
		metrics.RegisterPrometheusMetrics(),
		s.Certificate)
	if err != nil {
		log.Fatalf("could not connect to db: %q", err)
	}
	s.DBInterface = connectionInterface

	if err := connectionInterface.Ping(); err != nil {
		log.Fatalf("could not ping db %q", err)
	}

	// Generate seedyness
	rand.New(rand.NewSource(time.Now().UnixNano()))

	router := mux.NewRouter()
	router.Handle("/metrics", promhttp.Handler())
	work.Start(s)
	log.Fatal(http.ListenAndServe(":"+s.Port, router))

}
