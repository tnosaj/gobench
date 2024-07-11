package simple

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/tnosaj/gobench/internal"
)

// create basic prepare
// * tables
// * N rows
// * sk index

// Prepare stuff
func (st *SimpleReadWrite) Prepare() {
	logrus.Infof("prepare")

	err := st.S.DBInterface.AutoMigrateUP(fmt.Sprintf("%s/simple", st.S.SqlMigrationFolder))
	if err != nil {
		logrus.Errorf("Error when migrating: %q", err)
	}
	st.bulkInsert()
	logrus.Infof("Done")
}

func (st *SimpleReadWrite) bulkInsert() {
	wg := sync.WaitGroup{}
	ch := make(chan int)
	for i := 0; i < st.S.Concurrency; i++ {
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

func dbinsert(s *internal.Settings, tableName string) error {
	r := generateRow(s.Randomizer)
	err := s.DBInterface.ExecStatement("INSERT INTO "+tableName+"(k, c , pad) VALUES ("+strconv.Itoa(r.K)+",'"+r.C+"','"+r.Pad+"');", "blkinsert")
	if err != nil {
		logrus.Warnf("Error %s when inserting row into %s table. Values: %+v)", err, tableName, r)
		return err
	}
	return nil
}
