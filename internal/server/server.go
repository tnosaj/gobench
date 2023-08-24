package server

import (
	"github.com/sirupsen/logrus"
	"gitlab.otters.xyz/jason.tevnan/gobench/internal"
	"gitlab.otters.xyz/jason.tevnan/gobench/internal/strategy"
	"gitlab.otters.xyz/jason.tevnan/gobench/internal/strategy/insert"
	"gitlab.otters.xyz/jason.tevnan/gobench/internal/strategy/replica"
	"gitlab.otters.xyz/jason.tevnan/gobench/internal/strategy/simple"
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
