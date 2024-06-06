package server

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/tnosaj/gobench/internal"
	"github.com/tnosaj/gobench/internal/strategy"
	"github.com/tnosaj/gobench/internal/strategy/aerospike"
	"github.com/tnosaj/gobench/internal/strategy/insert"
	"github.com/tnosaj/gobench/internal/strategy/replica"
	"github.com/tnosaj/gobench/internal/strategy/simple"
)

type GobenchServer struct {
	Settings internal.Settings
	Strategy strategy.ExecutionStrategy
}

type HttpSettings struct {
	SqlMigrationFolder string `json:"sqlmigrationfolder"`
	Strategy           string `json:"strategy"`
	Concurrency        int    `json:"concurrency"`
	Initialdatasize    int    `json:"initialdatasize"`
	Duration           int    `json:"duration"`
	Rate               int    `json:"rate"`
}

func NewGobenchServer(settings internal.Settings) GobenchServer {
	logrus.Info("started server")

	var st strategy.ExecutionStrategy

	switch settings.Strategy {
	case "simple":
		st = simple.MakeSimpleStrategy(&settings)
	case "insert":
		st = insert.MakeInsertStrategy(&settings)
	case "replica":
		st = replica.MakeReplicaStrategy(&settings)
	case "aerospike":
		st = aerospike.MakeAerospikeStrategy(&settings)
	}
	var values []string

	if settings.TmpFile != "none" {
		if _, err := os.Stat(settings.TmpFile); errors.Is(err, os.ErrNotExist) {
			logrus.Debug("Tmpfile does not exist on startup")
		} else {
			logrus.Debugf("Tmpfile exists on startup, reading")
			values, err = readLines(settings.TmpFile)
			if err != nil {
				logrus.Errorf("There was a problem reading the tmpfile into memory: %s", err)
			}
		}
	}
	st.PopulateExistingValues(values)
	return GobenchServer{Settings: settings, Strategy: st}
}

func (gbs *GobenchServer) Shutdown(context context.Context) {
	logrus.Debug("Shutting down gobenchserver")
	list := gbs.Strategy.ReturnExistingValues()
	if gbs.Settings.TmpFile != "none" && len(list) > 0 {
		err := writeLines(list, gbs.Settings.TmpFile)

		if err != nil {
			logrus.Errorf("There was a problem writting the tmpfile onto disk: %s", err)
		}
	}
	gbs.Strategy.Shutdown(context)
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func writeLines(lines []string, path string) error {
	var file *os.File
	var err error

	if _, err = os.Stat(path); errors.Is(err, os.ErrNotExist) {
		logrus.Debug("Tmpfile does not exist on shutdown, creating")
		file, err = os.Create(path)
	} else {
		logrus.Debugf("Tmpfile exists on shutdown, appending")
		file, err = os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	}

	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}
