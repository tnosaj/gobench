package server

import (
	"github.com/sirupsen/logrus"
	"github.com/tnosaj/gobench/internal"
	"github.com/tnosaj/gobench/internal/strategy"
	"github.com/tnosaj/gobench/internal/strategy/insert"
	"github.com/tnosaj/gobench/internal/strategy/replica"
	"github.com/tnosaj/gobench/internal/strategy/simple"
)

type GobenchServer struct {
	Settings internal.Settings
	Strategy *strategy.ExecutionStrategy
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

	}
	return GobenchServer{Settings: settings, Strategy: &st}
}
