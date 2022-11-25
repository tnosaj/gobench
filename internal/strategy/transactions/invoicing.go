package transactions

import "github.com/sirupsen/logrus"

// invoicing_number_of_documents
// invoicing_document_created_datetime
// invoicing_document_number
// invoicing_document_status
// invoicing_exception_reason
// invoicing_document_attachment (pdf/html)
func (st TransactionReadWrite) createInvoicingEvent(txid string) {
	query := `INSERT INTO invoicingevents ()
VALUES ($1, $2, $3, $4)`
	var err error
	x := st.S.Randomizer.Intn(100)
	switch {
	case x <= 90:
		logrus.Debugf("No error in invoicing event")
		_, err = st.S.DBConnection.Exec(query, txid, "foo", "bar", "bas")
	default:
		logrus.Debugf("error in invoicing event")
		_, err = st.S.DBConnection.Exec(query, txid, "foo", "bar", nil)
	}
	if err != nil {
		logrus.Errorf("error running query %s", err)
	}
}
