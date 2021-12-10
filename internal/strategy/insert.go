package strategy

import (
	"strconv"

	"github.com/sirupsen/logrus"
	"gitlab.otters.xyz/jason.tevnan/gobench/internal"
	"gitlab.otters.xyz/jason.tevnan/gobench/internal/helper"
)

// InsertReadWrite do stuffs
type InsertReadWrite struct {
	S          internal.Settings
	MaxIDCount int
}

// CreateCommand do stuffs
func (st InsertReadWrite) CreateCommand() (string, string) {
	x := st.S.Randomizer.Intn(100)
	// x:50  - 50
	// r:100 - 0
	// w:0   - 100
	switch {
	case x <= st.S.ReadWriteSplit.Reads:
		logrus.Debugf("Will perform read")
		return st.read()
	default:
		logrus.Debugf("Will perform write")
		return st.write()
	}

}

func (st InsertReadWrite) read() (string, string) {
	switch st.S.Randomizer.Intn(3) {
	case 0, 1:
		logrus.Debugf("Will perform getPk")
		return st.getPk(), "getPk"
	default:
		logrus.Debugf("Will perform getSk")
		return st.getSk(), "getSk"
	}
}

func (st InsertReadWrite) write() (string, string) {
	logrus.Debugf("Will perform insert")
	return st.create(), "create"

}

// select by primary key
func (st InsertReadWrite) getPk() string {
	return "select id,k,c,pad from " + st.S.TableName + " where id=" + strconv.Itoa(st.S.Randomizer.Intn(st.MaxIDCount)) + ";"
}

// select by secondary key
func (st InsertReadWrite) getSk() string {
	return "select id,k,c,pad from " + st.S.TableName + " where k=" + strconv.Itoa(st.S.Randomizer.Intn(2147483647)) + ";"
}

// insert one record
func (st InsertReadWrite) create() string {
	r := helper.GenerateRow(st.S.Randomizer)
	return "INSERT INTO " + st.S.TableName + "(k, c , pad) VALUES (" + strconv.Itoa(r.K) + ",'" + r.C + "','" + r.Pad + "');"

}

// update one record
func (st InsertReadWrite) update() string {
	r := helper.GenerateRow(st.S.Randomizer)
	return "UPDATE " + st.S.TableName + " SET c='" + r.C + "', pad='" + r.Pad + "' WHERE id=" + strconv.Itoa(st.S.Randomizer.Intn(st.MaxIDCount)) + ";"
}

// delete one record
func (st InsertReadWrite) delete() string {
	return "DELETE FROM " + st.S.TableName + " WHERE id=" + strconv.Itoa(st.S.Randomizer.Intn(st.MaxIDCount)) + ";"
}
