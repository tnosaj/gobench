package insert

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/sirupsen/logrus"
	"gitlab.otters.xyz/jason.tevnan/gobench/internal"
	"gitlab.otters.xyz/jason.tevnan/gobench/internal/helper"
)

// create basic prepare
// * tables
// * N rows
// * sk index

// Prepare stuff
func (st InsertReadWrite) Prepare() {
	logrus.Infof("prepare")

	err := st.S.DBInterface.AutoMigrateUP(fmt.Sprintf("%s/insert", st.S.SqlMigrationFolder))
	if err != nil {
		logrus.Errorf("Error when migrating: %q", err)
	}
	st.bulkInsert()
	logrus.Infof("Done")
}

func (st InsertReadWrite) bulkInsert() {
	wg := sync.WaitGroup{}
	ch := make(chan int)
	for i := 0; i < 20; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			for range ch {
				err := dbinsert(st.S, st.TableName)
				if err != nil {
					logrus.Warnf("Error inserting: %s", err)
				}
			}
		}()
	}

	for i := 0; i < st.S.Initialdatasize; i++ {
		ch <- i
	}
	close(ch)
	wg.Wait()
}

func dbinsert(s internal.Settings, tableName string) error {
	r := helper.GenerateRow(s.Randomizer)
	err := s.DBInterface.ExecStatement("INSERT INTO "+tableName+"(k, c , pad) VALUES ("+strconv.Itoa(r.K)+",'"+r.C+"','"+r.Pad+"');", "blkinsert")
	if err != nil {
		logrus.Warnf("Error %s when inserting row into %s table. Values: %+v)", err, tableName, r)
		return err
	}
	return nil
}
