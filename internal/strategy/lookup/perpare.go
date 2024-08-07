package lookup

import (
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/tnosaj/gobench/internal"
)

// create basic prepare
// * tables
// * N rows
// * sk index

// Prepare stuff
func (a *Lookup) Prepare() {
	logrus.Infof("prepare")
	a.S.ServerStatus = "busy"
	err := a.S.DBInterface.AutoMigrateUP(fmt.Sprintf("%s/lookup", a.S.SqlMigrationFolder))
	if err != nil {
		logrus.Errorf("Error when migrating: %q", err)
	}
	a.bulkInsert()
	logrus.Infof("Done")
}

func (a *Lookup) bulkInsert() {
	wg := sync.WaitGroup{}
	ch := make(chan int)
	for i := 0; i < a.S.Concurrency; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			for range ch {
				err := a.dbinsert(a.S, a.StorageLocation)
				if err != nil {
					logrus.Warnf("Error inserting: %s", err)
				}
			}
		}()
	}

	for i := 0; i < a.S.Initialdatasize; i++ {
		ch <- i
	}
	close(ch)
	wg.Wait()
	a.S.ServerStatus = "free"
}

func (a *Lookup) dbinsert(s *internal.Settings, tableName string) error {
	r := a.create()
	err := s.DBInterface.ExecInterfaceStatement(r, "blkinsert")
	if err != nil {
		logrus.Warnf("Error %s when inserting row into %s table. Values: %+v)", err, tableName, r)
		return err
	}
	return nil
}
