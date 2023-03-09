package n26

import (
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
)

// create basic prepare
// * tables
// * N rows
// * sk index

// Prepare stuff
func (st N26ReadWrite) Prepare() {
	logrus.Infof("prepare")

	err := st.S.DBInterface.AutoMigrateUP(fmt.Sprintf("%s/n26", st.S.SqlMigrationFolder))
	if err != nil {
		logrus.Errorf("Error when migrating: %q", err)
	}
	st.bulkInsert()
	logrus.Infof("Done")
}

func (st N26ReadWrite) bulkInsert() {
	wg := sync.WaitGroup{}
	ch := make(chan int)
	for i := 0; i < 20; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			for range ch {
				err := st.dbinsert()
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

func (st N26ReadWrite) dbinsert() error {
	r := st.generateTransaction()
	err := st.S.DBInterface.ExecStatement(
		fmt.Sprintf("INSERT INTO transactions (id,account_id,card_id,amount,created) VALUES ('%s','%s','%s',%d,%d);",
			r.id, r.account_id, r.card_id, r.amount, r.created),
		"blkinsert",
	)
	if err != nil {
		logrus.Warnf("Error %s when inserting row into table.)", err)
		return err
	}
	return nil
}
