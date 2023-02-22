package insert

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

func (st InsertReadWrite) Cleanup() {
	logrus.Infof("cleanup")
	err := st.S.DBInterface.AutoMigrateDown(fmt.Sprintf("%s/insert", st.S.SqlMigrationFolder))
	if err != nil {
		logrus.Fatalf("Error when migrating: %s", err)
	}

	logrus.Infof("Done")
}
