package main

import (
	"math/rand"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/tnosaj/gobench/internal/db"
	"github.com/tnosaj/gobench/internal/server"
	"github.com/tnosaj/gobench/internal/startup"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	s, err := startup.EvaluateInputs()
	if err != nil {
		log.Fatalf("could not evaluate inputs: %q", err)
	}

	connectionInterface, connection, err := db.Connect(s.DB,
		s.DBConnectionInfo,
		s.Connectionpoolsize,
		s.TLSCerts)
	if err != nil {
		log.Fatalf("could not connect to db: %q", err)
	}
	s.DBConnection = connection
	s.DBInterface = connectionInterface

	if err := connectionInterface.Ping(); err != nil {
		log.Fatalf("could not ping db %q", err)
	}

	// Generate seedyness
	rand.New(rand.NewSource(time.Now().UnixNano()))

	server := server.NewGobenchServer(s)

	router := mux.NewRouter()
	router.Handle("/metrics", promhttp.Handler())
	router.HandleFunc("/status", server.Status).Methods("GET")
	router.HandleFunc("/prepare", server.Prepare).Methods("POST")
	router.HandleFunc("/run", server.Run).Methods("POST")
	router.HandleFunc("/cleanup", server.Cleanup).Methods("POST")
	log.Fatal(http.ListenAndServe(":"+s.Port, router))

}
