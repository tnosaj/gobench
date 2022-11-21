package transactions

import (
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gitlab.otters.xyz/jason.tevnan/gobench/internal"
)

// TransactionReadWrite do stuffs
type TransactionReadWrite struct {
	S          internal.Settings
	MaxIDCount int
}

func MakeTransactionReadWriteStrategy(s internal.Settings, action string) TransactionReadWrite {
	logrus.Info("creating TransactionReadWrite")
	if action == "run" {
		count, err := s.DBInterface.ExecStatementWithReturnInt("select count(id) from " + s.TableName + ";")

		if err != nil {
			logrus.Fatalf("could not get max id count with error: %q", err)
		}
		logrus.Infof("Query from 0 to %d", count)
		return TransactionReadWrite{
			S:          s,
			MaxIDCount: count,
		}
	}
	return TransactionReadWrite{
		S: s,
	}
}

// CreateCommand do stuffs
func (st TransactionReadWrite) RunCommand() {
	/*
		create reporting event in
		create transaction id
		* 1-10 orderlines
		* 1-3 tax per orderline
		create reporting event ack (10% failure)
		create invcoicing event in
		create invoicing event ack (10% failure)
		create tax event
		create lookup event
		20%(create returns event)
	*/
	transactionid := uuid.New().String()
	st.createReportingEvent(transactionid)
	st.createTransaction(transactionid)
	time.Sleep(100 * time.Millisecond)
	st.createReportingEvent(transactionid) // ack
	time.Sleep(100 * time.Millisecond)

	st.createInvoicingEvent(transactionid)
	time.Sleep(100 * time.Millisecond)
	st.createInvoicingEvent(transactionid) // ack
	time.Sleep(100 * time.Millisecond)

	st.createTaxEvent(transactionid)
	time.Sleep(100 * time.Millisecond)
	st.createTaxEvent(transactionid) // ack
	time.Sleep(100 * time.Millisecond)

	st.createLookupEvent(transactionid)
	time.Sleep(100 * time.Millisecond)
	st.createLookupEvent(transactionid) // ack
	time.Sleep(100 * time.Millisecond)

	st.createTaxEvent(transactionid)
	time.Sleep(100 * time.Millisecond)
	st.createTaxEvent(transactionid) // ack
	time.Sleep(100 * time.Millisecond)

}

func (st TransactionReadWrite) createTransaction(txid string) {
	query := `INSERT INTO transactions ()
		VALUES ($1, $2, $3, $4)`
	_, err := st.S.DBConnection.Exec(query, txid, "foo", "bar", "bas")
	if err != nil {
		logrus.Errorf("error running query %s", err)
	}
	for i := 0; i < st.S.Randomizer.Intn(10); i++ {
		st.createTransactionOrderLine(txid, i)
	}

}
func (st TransactionReadWrite) createTransactionOrderLine(txid string, orderlineid int) {
	query := `INSERT INTO transactions ()
	VALUES ($1, $2, $3, $4)`
	_, err := st.S.DBConnection.Exec(query, txid, orderlineid, "bar", "bas")
	if err != nil {
		logrus.Errorf("error running query %s", err)
	}
	for i := 0; i < st.S.Randomizer.Intn(3); i++ {
		st.createTransactionOrderLineTax(txid, orderlineid)
	}
}
func (st TransactionReadWrite) createTransactionOrderLineTax(txid string, orderlineid int) {
	query := `INSERT INTO transactions ()
	VALUES ($1, $2, $3, $4)`
	_, err := st.S.DBConnection.Exec(query, txid, orderlineid, st.getRandomCountry, "bas")
	if err != nil {
		logrus.Errorf("error running query %s", err)
	}
}

func (st TransactionReadWrite) getRandomCountry() string {
	countries := []string{"AT", "IT", "ES", "DE", "US", "GB", "PT", "HR", "HU", "PL", "NL", "FR"}
	return countries[st.S.Randomizer.Intn(len(countries))]
}

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
