package lookup

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
type Lookup struct {
	S               *internal.Settings
	MaxIDCount      int
	StorageLocation string
	Values          []string
}

func MakeLookupStrategy(s *internal.Settings) *Lookup {
	logrus.Info("creating Lookup")

	tableName := "sbtest"

	return &Lookup{
		S:               s,
		MaxIDCount:      s.Initialdatasize,
		StorageLocation: tableName,
	}
}

func (a *Lookup) PopulateExistingValues(v []string) {
	a.Values = v
}

func (a *Lookup) ReturnExistingValues() []string {
	return a.Values
}

func (a *Lookup) Shutdown(c context.Context) {
	logrus.Info("shutting down strategy")

	a.S.DBInterface.Shutdown(c)
}

func (a *Lookup) UpdateSettings(s internal.Settings) {
	a.S = &s
}

// CreateCommand do stuffs
func (a *Lookup) RunCommand() {
	x := a.S.Randomizer.Intn(100)
	// x:50  - 50
	// r:100 - 0
	// w:0   - 100
	switch {
	case x <= a.S.ReadWriteSplit.Reads:
		logrus.Debugf("Will perform read")
		a.S.DBInterface.ExecInterfaceStatement(a.read())
	default:
		logrus.Debugf("Will perform write")
		a.S.DBInterface.ExecInterfaceStatement(a.write())
	}

}

func (a *Lookup) read() (string, string) {
	logrus.Debugf("Will perform getRandom")
	return a.getRandom(), "read"

}

func (a *Lookup) write() (string, string) {
	logrus.Debugf("Will perform insert")
	return a.create(), "write"

}

// select by primary key
func (a *Lookup) getRandom() string {
	var b bytes.Buffer
	b.WriteString(a.StorageLocation)
	b.WriteString(",")
	b.WriteString(a.randomUUIDFromList(a.S.Randomizer))
	return b.String()
}

// insert one record
func (a *Lookup) create() string {
	id := uuid.New().String()
	var b bytes.Buffer
	b.WriteString(a.StorageLocation)
	b.WriteString(",")
	b.WriteString(id)
	b.WriteString(",")
	b.WriteString(id)
	a.Values = append(a.Values, id)
	return b.String()

}

func (a *Lookup) randomUUIDFromList(rand internal.Random) string {
	randomIndex := rand.Intn(len(a.Values))
	return a.Values[randomIndex]

}
