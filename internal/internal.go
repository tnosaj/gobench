package internal

import (
	"gitlab.otters.xyz/jason.tevnan/gobench/internal/db"
	"gitlab.otters.xyz/jason.tevnan/gobench/pkg/args"
)

// Settings contains global settings
type Settings struct {
	Concurrency        int
	Connectionpoolsize int
	Initialdatasize    int
	Duration           int
	Rate               int

	Debug bool

	Action           string
	DB               string
	DBConnectionInfo db.ConnectionInfo
	DurationType     string
	Port             string
	TableName        string

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
