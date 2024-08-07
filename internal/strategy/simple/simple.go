package simple

import (
	"context"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/tnosaj/gobench/internal"
)

// Row of work
type Row struct {
	ID  int
	K   int
	C   string
	Pad string
}

// SimpleReadWrite do stuffs
type SimpleReadWrite struct {
	S          *internal.Settings
	MaxIDCount int
	TableName  string
	Values     []string
}

func MakeSimpleStrategy(s *internal.Settings) *SimpleReadWrite {
	logrus.Info("creating SimpleReadWrite")

	tableName := "sbtest"

	// if action == "run" {
	// 	count, err := s.DBInterface.ExecStatementWithReturnInt("select count(id) from " + tableName + ";")

	// 	if err != nil {
	// 		logrus.Fatalf("could not get max id count with error: %q", err)
	// 	}
	// 	logrus.Infof("Query from 0 to %d", count)
	// 	return &SimpleReadWrite{
	// 		S:          s,
	// 		MaxIDCount: count,
	// 		TableName:  tableName,
	// 	}
	// }
	return &SimpleReadWrite{
		S:          s,
		MaxIDCount: s.Initialdatasize,
		TableName:  tableName,
	}
}

func (st *SimpleReadWrite) UpdateSettings(s *internal.Settings) {
	st.S = s
}

func (st *SimpleReadWrite) PopulateExistingValues(v []string) {
	st.Values = v
}

func (st *SimpleReadWrite) ReturnExistingValues() []string {
	return st.Values
}

func (st *SimpleReadWrite) Shutdown(c context.Context) {
	st.Shutdown(c)
}

// CreateCommand do stuffs
func (st *SimpleReadWrite) RunCommand() {
	x := st.S.Randomizer.Intn(100)
	// x:50  - 50
	// r:100 - 0
	// w:0   - 100
	switch {
	case x <= st.S.ReadWriteSplit.Reads:
		logrus.Debugf("Will perform read")
		st.S.DBInterface.ExecStatement(st.read())
	default:
		logrus.Debugf("Will perform write")
		st.S.DBInterface.ExecStatement(st.write())
	}

}

func (st *SimpleReadWrite) read() (string, string) {
	switch st.S.Randomizer.Intn(3) {
	case 0, 1:
		logrus.Debugf("Will perform getPk")
		return st.getPk(), "getPk"
	default:
		logrus.Debugf("Will perform getSk")
		return st.getSk(), "getSk"
	}
}

func (st *SimpleReadWrite) write() (string, string) {
	switch st.S.Randomizer.Intn(3) {
	case 0:
		logrus.Debugf("Will perform insert")
		return st.create(), "create"
	case 1:
		logrus.Debugf("Will perform delete")
		return st.delete(), "delete"
	default:
		logrus.Debugf("Will perform update")
		return st.update(), "update"
	}
}

// select by primary key
func (st *SimpleReadWrite) getPk() string {
	return "select id,k,c,pad from " + st.TableName + " where id=" + strconv.Itoa(st.S.Randomizer.Intn(st.MaxIDCount)) + ";"
}

// select by secondary key
func (st *SimpleReadWrite) getSk() string {
	return "select id,k,c,pad from " + st.TableName + " where k=" + strconv.Itoa(st.S.Randomizer.Intn(2147483647)) + ";"
}

// insert one record
func (st *SimpleReadWrite) create() string {
	r := generateRow(st.S.Randomizer)
	return "INSERT INTO " + st.TableName + "(k, c , pad) VALUES (" + strconv.Itoa(r.K) + ",'" + r.C + "','" + r.Pad + "');"

}

// update one record
func (st *SimpleReadWrite) update() string {
	r := generateRow(st.S.Randomizer)
	return "UPDATE " + st.TableName + " SET c='" + r.C + "', pad='" + r.Pad + "' WHERE id=" + strconv.Itoa(st.S.Randomizer.Intn(st.MaxIDCount)) + ";"
}

// delete one record
func (st *SimpleReadWrite) delete() string {
	return "DELETE FROM " + st.TableName + " WHERE id=" + strconv.Itoa(st.S.Randomizer.Intn(st.MaxIDCount)) + ";"
}

// generateRow returns a row
func generateRow(rand internal.Random) Row {
	return Row{
		K:   rand.Intn(2147483647),
		C:   randomString(120, rand),
		Pad: randomString(60, rand),
	}
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randomString(n int, rand internal.Random) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
