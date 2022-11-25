package transactions

import (
	"github.com/sirupsen/logrus"
)

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
