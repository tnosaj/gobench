package lookup

import (
	"bytes"
	"context"

	"github.com/samborkent/uuid"
	"github.com/sirupsen/logrus"
	"github.com/tnosaj/gobench/internal"
	"github.com/tnosaj/gobench/internal/cache"
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
	ValuesInterface cache.CacheValues
}

func MakeLookupStrategy(s *internal.Settings, cache cache.CacheValues) *Lookup {
	logrus.Info("creating Lookup")

	tableName := "sbtest"

	return &Lookup{
		S:               s,
		MaxIDCount:      s.Initialdatasize,
		StorageLocation: tableName,
		ValuesInterface: cache,
	}
}

func (a *Lookup) Shutdown(c context.Context) {
	logrus.Info("shutting down strategy")
	a.ValuesInterface.Save()
	a.S.DBInterface.Shutdown(c)
}

func (a *Lookup) UpdateSettings(s *internal.Settings) {
	a.S = s
}

// CreateCommand do stuffs
func (a *Lookup) RunCommand() {
	x := a.S.Randomizer.Intn(100)
	// x:50  - 50
	// r:100 - 0
	// w:0   - 100
	switch {
	case x <= a.S.ReadWriteSplit.Reads:
		logrus.Tracef("Will perform read")
		a.S.DBInterface.ExecInterfaceStatement(a.read())
	default:
		logrus.Tracef("Will perform write")
		a.S.DBInterface.ExecInterfaceStatement(a.write())
	}

}

func (a *Lookup) read() (string, string) {
	logrus.Tracef("Will perform getRandom")
	// 50:50::hit:miss ratio
	if a.S.Randomizer.Intn(100) <= 50 {
		return a.getRandom(), "read"
	}
	return a.getFailingRandom(), "read-404"

}

func (a *Lookup) write() (string, string) {
	logrus.Tracef("Will perform insert")
	return a.create(), "write"

}

// select by primary key from existing list
func (a *Lookup) getRandom() string {
	var b bytes.Buffer
	uid, err := a.ValuesInterface.GetRandom()
	if err != nil {
		return a.getFailingRandom() // then we will just fail, thats fine
	}
	b.WriteString(a.StorageLocation)
	b.WriteString(",")
	b.WriteString(uid.String())
	return b.String()
}

// failing select by pk
func (a *Lookup) getFailingRandom() string {
	var b bytes.Buffer
	b.WriteString(a.StorageLocation)
	b.WriteString(",")
	b.WriteString(uuid.New().String())
	return b.String()
}

// insert one record
func (a *Lookup) create() string {
	id := uuid.New()
	var b bytes.Buffer
	b.WriteString(a.StorageLocation)
	b.WriteString(",")
	b.WriteString(id.String())
	b.WriteString(",")
	b.WriteString(id.String())
	err := a.ValuesInterface.Put(id)
	if err != nil {
		logrus.Errorf("Error adding value to cache: %s", err)
	}
	return b.String()

}
