package replica

import (
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
)

// create basic prepare
// * tables
// * N rows
// * sk index

// 20 regions
// 100k rows

// Prepare stuff
func (st *ReplicaReadWrite) Prepare() {
	logrus.Infof("prepare")

	err := st.S.DBInterface.AutoMigrateUP(fmt.Sprintf("%s/replica", st.S.SqlMigrationFolder))
	if err != nil {
		logrus.Errorf("Error when migrating: %q", err)
	}
	st.bulkInsert()
	logrus.Debugf("length: %d", len(st.Assets))
	logrus.Infof("Done")
}

func (st *ReplicaReadWrite) bulkInsert() {
	wg := sync.WaitGroup{}
	ch := make(chan int)
	for i := 0; i < 20; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			for range ch {
				err := st.S.DBInterface.ExecStatement(st.create(), "blikinsert")
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
