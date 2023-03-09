package n26

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gitlab.otters.xyz/jason.tevnan/gobench/internal"
)

type transaction struct {
	id         string
	card_id    string
	account_id string
	amount     int
	created    int64
}

// N26ReadWrite do stuffs
type N26ReadWrite struct {
	S          internal.Settings
	MaxIDCount int
	Accounts   []string
	Cards      []string
}

func MakeN26Strategy(s internal.Settings, action string) N26ReadWrite {
	logrus.Info("creating N26ReadWrite")

	accounts := []string{
		"46f5d9a5-6292-4b67-a32b-90e44edea5c1",
		"1e3b2691-43bb-4825-b0b1-85380be9512d",
		"efd50ed6-0bc3-44a9-881f-ec89713fdd80",
	}
	cards := []string{
		"ee7449d0-e8e6-464c-9989-aefd995fb50d",
		"26501dfa-e783-4443-a70a-8a3a5f8877f0",
		"69b73b10-581d-4912-9833-6d5cd9cfc3d7",
		"e380f48d-fe4f-4946-8286-b48cf9967aa5",
	}

	return N26ReadWrite{
		S:        s,
		Accounts: accounts,
		Cards:    cards,
	}
}

// CreateCommand do stuffs
func (st N26ReadWrite) RunCommand() {
	logrus.Debugf("Will perform write")
	st.S.DBInterface.ExecStatement(st.write())

}

func (st N26ReadWrite) write() (string, string) {
	logrus.Debugf("Will perform insert")
	return st.create(), "create"
}

// insert one record
func (st N26ReadWrite) create() string {
	r := st.generateTransaction()
	return fmt.Sprintf("INSERT INTO transactions (id,account_id,card_id,amount,created) VALUES ('%s','%s','%s',%d,%d);",
		r.id, r.account_id, r.card_id, r.amount, r.created,
	)
}

// generateTransaction returns a row
func (st N26ReadWrite) generateTransaction() transaction {
	return transaction{
		id:         uuid.New().String(),
		account_id: st.Accounts[st.S.Randomizer.Intn(len(st.Accounts))],
		card_id:    st.Cards[st.S.Randomizer.Intn(len(st.Cards))],
		amount:     st.S.Randomizer.Intn(3000),
		created:    time.Now().UnixNano() / int64(time.Millisecond),
	}
}
