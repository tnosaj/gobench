package transactions

import (
	"github.com/sirupsen/logrus"
)

func (st TransactionReadWrite) Cleanup() {
	logrus.Infof("cleanup")
	err := st.automigrate("cleanup")
	if err != nil {
		logrus.Fatalf("Could not cleanup with error %s", err)
	}
	logrus.Infof("Done, please end with ctl+c")
}
