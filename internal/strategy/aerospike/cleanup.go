package aerospike

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

func (a Aerospike) Cleanup() {
	logrus.Infof("cleanup")
	err := a.S.DBInterface.AutoMigrateDown(fmt.Sprintf("%s/simple", a.S.SqlMigrationFolder))
	if err != nil {
		logrus.Fatalf("Error when migrating: %s", err)
	}

	logrus.Infof("Done")
}
