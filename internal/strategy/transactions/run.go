package transactions

import (
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/tnosaj/gobench/internal"
)

// TransactionReadWrite do stuffs
type TransactionReadWrite struct {
	S internal.Settings
}

func MakeTransactionReadWriteStrategy(s internal.Settings, action string) TransactionReadWrite {
	logrus.Info("creating TransactionReadWrite")
	return TransactionReadWrite{
		S: s,
	}
}

// CreateCommand do stuffs
func (st TransactionReadWrite) RunCommand() {

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
