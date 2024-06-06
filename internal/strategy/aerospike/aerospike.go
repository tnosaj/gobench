package aerospike

import (
	"bytes"
	"context"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/tnosaj/gobench/internal"
)

// Row of work
type Row struct {
	ID  int
	K   int
	C   string
	Pad string
}

// SimpleReadWrite do stuffs
type Aerospike struct {
	S          *internal.Settings
	MaxIDCount int
	NameSpace  string
	SetName    string
	Values     []string
}

func MakeAerospikeStrategy(s *internal.Settings) *Aerospike {
	logrus.Info("creating Aerospike")

	tableName := "sbtest"

	return &Aerospike{
		S:          s,
		MaxIDCount: s.Initialdatasize,
		NameSpace:  tableName,
		SetName:    tableName,
	}
}

func (a *Aerospike) PopulateExistingValues(v []string) {
	a.Values = v
}

func (a *Aerospike) ReturnExistingValues() []string {
	return a.Values
}

func (a *Aerospike) Shutdown(c context.Context) {
	logrus.Info("shutting down strategy")

	a.S.DBInterface.Shutdown(c)
}

func (a *Aerospike) UpdateSettings(s internal.Settings) {
	a.S = &s
}

// CreateCommand do stuffs
func (a *Aerospike) RunCommand() {
	x := a.S.Randomizer.Intn(100)
	// x:50  - 50
	// r:100 - 0
	// w:0   - 100
	switch {
	case x <= a.S.ReadWriteSplit.Reads:
		logrus.Debugf("Will perform read")
		a.S.DBInterface.ExecStatement(a.read())
	default:
		logrus.Debugf("Will perform write")
		a.S.DBInterface.ExecStatement(a.write())
	}

}

func (a *Aerospike) read() (string, string) {
	logrus.Debugf("Will perform getRandom")
	return a.getRandom(), "read"

}

func (a *Aerospike) write() (string, string) {
	logrus.Debugf("Will perform insert")
	return a.create(), "write"

}

// select by primary key
func (a *Aerospike) getRandom() string {
	var b bytes.Buffer
	b.WriteString(a.NameSpace)
	b.WriteString(",")
	b.WriteString(a.SetName)
	b.WriteString(",")
	b.WriteString(a.randomUUIDFromList(a.S.Randomizer))
	return b.String()
}

// insert one record
func (a *Aerospike) create() string {
	id := uuid.New().String()
	var b bytes.Buffer
	b.WriteString(a.NameSpace)
	b.WriteString(",")
	b.WriteString(a.SetName)
	b.WriteString(",")
	b.WriteString(id)
	b.WriteString(",")
	b.WriteString(id)
	a.Values = append(a.Values, id)
	return b.String()

}

func (a *Aerospike) randomUUIDFromList(rand internal.Random) string {
	randomIndex := rand.Intn(len(a.Values))
	return a.Values[randomIndex]

}
