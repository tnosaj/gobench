package internal

import (
	"database/sql"
	"fmt"

	"github.com/tnosaj/gobench/internal/db"
	"github.com/tnosaj/gobench/pkg/args"
)

// Settings contains global settings
type Settings struct {
	Concurrency        int
	Connectionpoolsize int
	Initialdatasize    int
	Duration           int
	Rate               int

	Debug bool

	Action             string
	DB                 string
	DBConnection       *sql.DB
	DBConnectionInfo   db.ConnectionInfo
	DurationType       string
	Port               string
	SqlMigrationFolder string

	Strategy string

	TLSCerts db.TLSCerts

	DBInterface    db.DB
	ReadWriteSplit args.ReadWriteSplit
	Randomizer     Random
}

// Random interface to help with testing
type Random interface {
	Intn(int) int
}

func (settingOriginal Settings) PrintableSettings() string {
	type fakeSetting Settings
	printable := fakeSetting(settingOriginal)
	printable.DBConnectionInfo.Password = "[REDACTED]"
	return fmt.Sprintf("%#v", printable)
}
