package simple

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

func (st SimpleReadWrite) Cleanup() {
	logrus.Infof("cleanup")
	err := st.S.DBInterface.AutoMigrateDown(fmt.Sprintf("%s/simple", st.S.SqlMigrationFolder))
	if err != nil {
		logrus.Fatalf("Error when migrating: %s", err)
	}

	logrus.Infof("Done")
}
