package transactions

import "github.com/sirupsen/logrus"

// tax_calculation_request_datetime
// tax_logic_label
func (st TransactionReadWrite) createTaxEvent(txid string) {
	query := `INSERT INTO taxevents ()
VALUES ($1, $2, $3, $4)`
	var err error
	x := st.S.Randomizer.Intn(100)
	switch {
	case x <= 90:
		logrus.Debugf("No error in tax event")
		_, err = st.S.DBConnection.Exec(query, txid, "foo", "bar", "bas")
	default:
		logrus.Debugf("error in tax event")
		_, err = st.S.DBConnection.Exec(query, txid, "foo", "bar", nil)
	}
	if err != nil {
		logrus.Errorf("error running query %s", err)
	}
}
