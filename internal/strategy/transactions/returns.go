package transactions

import "github.com/sirupsen/logrus"

// returns_filing_id
// returns_filing_status
// returns_filing_deadline
func (st TransactionReadWrite) createReturnsEvent(txid string) {
	query := `INSERT INTO returnsevents ()
	VALUES ($1, $2, $3, $4)`
	var err error
	x := st.S.Randomizer.Intn(100)
	switch {
	case x <= 90:
		logrus.Debugf("No error in returns event")
		_, err = st.S.DBConnection.Exec(query, txid, "foo", "bar", "bas")
	default:
		logrus.Debugf("error in returns event")
		_, err = st.S.DBConnection.Exec(query, txid, "foo", "bar", nil)
	}
	if err != nil {
		logrus.Errorf("error running query %s", err)
	}
}
