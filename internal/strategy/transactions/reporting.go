package transactions

import "github.com/sirupsen/logrus"

// reporting_transaction_reported_datetime
// reporting_status
// reporting_reference_id
// reporting_job_id
// reporting_error_count
// reporting_error_datetime
// reporting_error_type
// reporting_error_code
// reporting_error_message
// reporting_error_field
// reporting_error_status
// reporting_request_json
func (st TransactionReadWrite) createReportingEvent(txid string) {
	query := `INSERT INTO reportingevents ()
		VALUES ($1, $2, $3, $4)`
	var err error
	x := st.S.Randomizer.Intn(100)
	switch {
	case x <= 90:
		logrus.Debugf("No error in reporting event")
		_, err = st.S.DBConnection.Exec(query, txid, "foo", "bar", "bas")
	default:
		logrus.Debugf("error in reporting event")
		_, err = st.S.DBConnection.Exec(query, txid, "foo", "bar", nil)
	}
	if err != nil {
		logrus.Errorf("error running query %s", err)
	}
}
