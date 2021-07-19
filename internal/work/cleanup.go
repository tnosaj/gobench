package work

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gitlab.otters.xyz/jason.tevnan/gobench/internal"
)

func cleanup(s internal.Settings, wp *workerPool) {
	logrus.Infof("cleanup")
	err := dropTable(s)
	if err != nil {
		logrus.Fatalf("Error when dropping table: %s", err)
	}

	logrus.Infof("Done, please end with ctl+c")
}

func dropTable(s internal.Settings) error {
	checkStatement := s.DBInterface.GetTableExists(s.DBConnectionInfo.DBName, s.TableName)
	exists, err := s.DBInterface.ExecStatementWithReturnBool(checkStatement)

	if err != nil {
		return fmt.Errorf("could not check if table '%s.%s' exists with error: %q", s.DBConnectionInfo.DBName, s.TableName, err)
	}
	if !exists {
		logrus.Infof("Nothing to drop, done")
		return nil
	}
	statement := `DROP TABLE ` + s.TableName + `;`
	err = s.DBInterface.ExecDDLStatement(statement)
	if err != nil {
		return fmt.Errorf("Error %s when dropping %s.%s table", err, s.DBConnectionInfo.DBName, s.TableName)
	}
	return nil
}
