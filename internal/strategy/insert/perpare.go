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

	err := createTable(st.S)
	if err != nil {
		logrus.Errorf("Error when creating table '%s.%s': %q", st.S.DBConnectionInfo.DBName, st.S.TableName, err)
	}
	bulkInsert(st.S)
	logrus.Infof("Done, please end with ctl+c")
}

func createTable(s internal.Settings) error {
	statement := s.DBInterface.GetTableExists(s.DBConnectionInfo.DBName, s.TableName)
	exists, err := s.DBInterface.ExecStatementWithReturnBool(statement)

	if err != nil {
		return fmt.Errorf("could not check if table '%s.%s' exists with error: %q", s.DBConnectionInfo.DBName, s.TableName, err)
	}
	if exists {
		return nil
	}
	query := s.DBInterface.Createable(s.DBConnectionInfo.DBName, s.TableName)
	err = s.DBInterface.ExecDDLStatement(query)
	if err != nil {
		return fmt.Errorf("could not create table '%s.%s' with error: %q", s.DBConnectionInfo.DBName, s.TableName, err)
	}
	return nil
}

func bulkInsert(s internal.Settings) {
	wg := sync.WaitGroup{}
	ch := make(chan int)
	for i := 0; i < 20; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			for range ch {
				err := dbinsert(s)
				if err != nil {
					logrus.Warnf("Error inserting: %s", err)
				}
			}
		}()
	}

	for i := 0; i < s.Initialdatasize; i++ {
		ch <- i
	}
	close(ch)
	wg.Wait()
}

func dbinsert(s internal.Settings) error {
	r := helper.GenerateRow(s.Randomizer)
	err := s.DBInterface.ExecStatement("INSERT INTO "+s.TableName+"(k, c , pad) VALUES ("+strconv.Itoa(r.K)+",'"+r.C+"','"+r.Pad+"');", "blkinsert")
	if err != nil {
		logrus.Warnf("Error %s when inserting row into %s table. Values: %+v)", err, s.TableName, r)
		return err
	}
	return nil
}
