package n26

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

func (st N26ReadWrite) Cleanup() {
	logrus.Infof("cleanup")
	err := st.S.DBInterface.AutoMigrateDown(fmt.Sprintf("%s/n26", st.S.SqlMigrationFolder))
	if err != nil {
		logrus.Fatalf("Error when migrating: %s", err)
	}

	logrus.Infof("Done")
}
