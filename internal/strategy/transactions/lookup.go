package transactions

import "github.com/sirupsen/logrus"

// customer_tin_validation_outcome
// customer_tin_validation_error_reason
// customer_tin_validation_request_datetime
func (st TransactionReadWrite) createLookupEvent(txid string) {
	query := `INSERT INTO lookupevents ()
	VALUES ($1, $2, $3, $4)`
	var err error
	x := st.S.Randomizer.Intn(100)
	switch {
	case x <= 90:
		logrus.Debugf("No error in lookup event")
		_, err = st.S.DBConnection.Exec(query, txid, "foo", "bar", "bas")
	default:
		logrus.Debugf("error in lookup event")
		_, err = st.S.DBConnection.Exec(query, txid, "foo", "bar", nil)
	}
	if err != nil {
		logrus.Errorf("error running query %s", err)
	}
}
