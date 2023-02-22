package server

import (
	"github.com/sirupsen/logrus"
	"gitlab.otters.xyz/jason.tevnan/gobench/internal"
)

type GobenchServer struct {
	Settings internal.Settings
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
	return GobenchServer{Settings: settings}
}
