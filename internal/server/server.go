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
	"github.com/tnosaj/gobench/internal/strategy/insert"
	"github.com/tnosaj/gobench/internal/strategy/lookup"
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
	logrus.Info("starting server")

	var st strategy.ExecutionStrategy

	switch settings.Strategy {
	case "simple":
		st = simple.MakeSimpleStrategy(&settings)
	case "insert":
		st = insert.MakeInsertStrategy(&settings)
	case "replica":
		st = replica.MakeReplicaStrategy(&settings)
	case "lookup":
		st = lookup.MakeLookupStrategy(&settings)
	default:
		logrus.Fatalf("unknown strategy: %s", settings.Strategy)
	}
	var values []string

	if settings.TmpFile != "none" {
		if _, err := os.Stat(fmt.Sprintf("%s-bac", settings.TmpFile)); err == nil {
			logrus.Fatalf("Unclean shutdown, tmpfile still exists: %s", fmt.Sprintf("%s-bac", settings.TmpFile))
		}
		if _, err := os.Stat(settings.TmpFile); errors.Is(err, os.ErrNotExist) {
			logrus.Info("Tmpfile does not exist on startup")
		} else {
			logrus.Info("Tmpfile exists on startup, reading")
			values, err = readLines(settings.TmpFile)
			if err != nil {
				logrus.Errorf("There was a problem reading the tmpfile into memory: %s", err)
			}
			err = os.Rename(settings.TmpFile, fmt.Sprintf("%s-bac", settings.TmpFile))
			if err != nil {
				logrus.Errorf("There was a problem moving the tmpfile to backup: %s", err)
			}
			logrus.Info("finished reading tmpfile")
		}
	}
	st.PopulateExistingValues(values)
	settings.ServerStatus = "free"
	logrus.Info("startup server complete")
	return GobenchServer{Settings: settings, Strategy: st}
}

func (gbs *GobenchServer) Shutdown(context context.Context) {
	logrus.Debug("Shutting down gobenchserver")
	list := gbs.Strategy.ReturnExistingValues()
	if gbs.Settings.TmpFile != "none" && len(list) > 0 {
		logrus.Info("start writting tmp file")
		err := writeLines(list, gbs.Settings.TmpFile)
		if err != nil {
			logrus.Errorf("There was a problem writting the tmpfile onto disk: %s", err)
		}
		logrus.Info("end writting tmp file")
		e := os.Remove(fmt.Sprintf("%s-bac", gbs.Settings.TmpFile))
		if e != nil {
			logrus.Errorf("There was a problem removing backup tmpfile from disk: %s", err)
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
