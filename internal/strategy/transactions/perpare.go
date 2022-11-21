package transactions

import (
	"github.com/sirupsen/logrus"
)

// create basic prepare
// * tables
// * N rows
// * sk index

// Prepare stuff
func (st TransactionReadWrite) Prepare() {
	logrus.Infof("prepare")
	err := st.automigrate("prepare")
	if err != nil {
		logrus.Fatalf("Could not prepare with error %s", err)
	}
	logrus.Infof("Done, please end with ctl+c")
}
