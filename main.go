package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/tnosaj/gobench/internal/db"
	"github.com/tnosaj/gobench/internal/server"
	"github.com/tnosaj/gobench/internal/startup"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	var wait time.Duration
	wait = time.Second * time.Duration(10)

	s, err := startup.EvaluateInputs()
	if err != nil {
		log.Fatalf("could not evaluate inputs: %q", err)
	}

	connectionInterface, err := db.Connect(s.DB,
		s.DBConnectionInfo,
		s.Connectionpoolsize,
		s.TLSCerts)
	if err != nil {
		log.Fatalf("could not connect to db: %q", err)
	}
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

	srv := &http.Server{
		Addr: fmt.Sprintf("0.0.0.0:%s", s.Port),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: wait,
		ReadTimeout:  wait,
		IdleTimeout:  wait,
		Handler:      router,
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logrus.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	server.Shutdown(ctx)
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)

}
