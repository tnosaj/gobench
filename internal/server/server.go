package server

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/tnosaj/gobench/internal"
	"github.com/tnosaj/gobench/internal/cache"
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
	Split              string `json:"split"`
	Concurrency        int    `json:"concurrency"`
	Initialdatasize    int    `json:"initialdatasize"`
	Duration           int    `json:"duration"`
	Rate               int    `json:"rate"`
}

func NewGobenchServer(settings internal.Settings) GobenchServer {
	logrus.Info("starting server")

	cache := cache.NewCache(settings.CacheType, settings.Randomizer)
	if err := cache.Load(); err != nil {
		logrus.Errorf("Could not load cache, continue without cache")
	}
	var st strategy.ExecutionStrategy

	switch settings.Strategy {
	case "simple":
		st = simple.MakeSimpleStrategy(&settings)
	case "insert":
		st = insert.MakeInsertStrategy(&settings)
	case "replica":
		st = replica.MakeReplicaStrategy(&settings)
	case "lookup":
		st = lookup.MakeLookupStrategy(&settings, cache)

	default:
		logrus.Fatalf("unknown strategy: %s", settings.Strategy)
	}

	settings.ServerStatus = "free"
	logrus.Info("startup server complete")
	return GobenchServer{Settings: settings, Strategy: st}
}

func (gbs *GobenchServer) Shutdown(context context.Context) {
	logrus.Debug("Shutting down gobenchserver")
	gbs.Strategy.Shutdown(context)
}
