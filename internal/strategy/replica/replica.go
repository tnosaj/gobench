package replica

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/tnosaj/gobench/internal"
)

// Row of work
type Row struct {
	ID                 string
	Driver_id          string
	Connected          bool
	Created_at         time.Time
	Loc                [2]float64
	Bearing            float32
	Altitude           float32
	Region_id          string
	Queue_zone_id      string
	Queue_zone_left_id string
	Queue_zone_left_at time.Time
	Queue_zone_hit_at  time.Time
	Queue_zone_loc     [2]float64
	Loc_ts             time.Time
	Speed              float32
	Alu_serial         int64
	Snoozed_at         time.Time
}

// ReplicaReadWrite do stuffs
type ReplicaReadWrite struct {
	S         *internal.Settings
	TableName string
	Regions   [20]string
	Assets    []string
}

func MakeReplicaStrategy(s *internal.Settings) *ReplicaReadWrite {
	logrus.Info("creating ReplicaReadWrite")

	tableName := "sbtest"

	var regions [20]string
	for i := 0; i < 20; i++ {
		regions[i] = uuid.NewString()
	}

	return &ReplicaReadWrite{
		S:         s,
		TableName: tableName,
		Regions:   regions,
	}
}

func (st *ReplicaReadWrite) UpdateSettings(s internal.Settings) {
	st.S = &s
}

// CreateCommand do stuffs
func (st *ReplicaReadWrite) RunCommand() {
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

func (st *ReplicaReadWrite) read() (string, string) {
	switch st.S.Randomizer.Intn(3) {
	case 0, 1:
		logrus.Debugf("Will perform getPk")
		return st.getPk(), "getPk"
	default:
		logrus.Debugf("Will perform getSk")
		return st.getSk(), "getSk"
	}
}

func (st *ReplicaReadWrite) write() (string, string) {

	switch st.S.Randomizer.Intn(5) {
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
func (st *ReplicaReadWrite) getPk() string {
	id := st.getRandomAsset()
	return fmt.Sprintf("select id, driver_id, connected, created_at, loc, bearing, altitude, region_id, queue_zone_id, queue_zone_left_id, queue_zone_left_at, queue_zone_hit_at, queue_zone_loc, loc_ts, speed, alu_serial, snoozed_at from %s where region_id='%s';", st.TableName, id)
}

// select by secondary key
func (st *ReplicaReadWrite) getSk() string {
	region_id := st.getRandomRegion()
	return fmt.Sprintf("select id, driver_id, connected, created_at, loc, bearing, altitude, region_id, queue_zone_id, queue_zone_left_id, queue_zone_left_at, queue_zone_hit_at, queue_zone_loc, loc_ts, speed, alu_serial, snoozed_at from %s where region_id='%s';", st.TableName, region_id)
}

// insert one record
func (st *ReplicaReadWrite) create() string {
	r := generateRow(st.S.Randomizer)
	r.Region_id = st.getRandomRegion()
	st.Assets = append(st.Assets, r.ID)
	return fmt.Sprintf("INSERT INTO %s (id, driver_id, connected, created_at, loc, bearing, altitude, region_id, queue_zone_id, queue_zone_left_id, queue_zone_left_at, queue_zone_hit_at, queue_zone_loc, loc_ts, speed, alu_serial, snoozed_at) VALUES ('%s','%s',%t,'%s',ST_GeomFromText('POINT(%f %f)'),%f,%f,'%s','%s','%s','%s','%s',ST_GeomFromText('POINT(%f %f)'),'%s',%f,%d,'%s');",
		st.TableName,
		r.ID, r.Driver_id, r.Connected, r.Created_at.Format("2006-01-02 15:04:05"),
		r.Loc[0], r.Loc[1], r.Bearing, r.Altitude, r.Region_id,
		r.Queue_zone_id, r.Queue_zone_left_id, r.Queue_zone_left_at.Format("2006-01-02 15:04:05"), r.Queue_zone_hit_at.Format("2006-01-02 15:04:05"),
		r.Queue_zone_loc[0], r.Queue_zone_loc[1], r.Loc_ts.Format("2006-01-02 15:04:05"), r.Speed, r.Alu_serial, r.Snoozed_at.Format("2006-01-02 15:04:05"),
	)
}

// update one record
func (st *ReplicaReadWrite) update() string {
	r := generateRow(st.S.Randomizer)
	r.Region_id = st.getRandomRegion()
	r.ID = st.getRandomAsset()
	return fmt.Sprintf("UPDATE %s SET bearing=%f, loc=ST_GeomFromText('POINT(%f %f)'), loc_ts='%s' WHERE id='%s';", st.TableName, r.Bearing, r.Loc[0], r.Loc[1], r.ID, time.Now().Format("2006-01-02 15:04:05"))
}

// delete one record
func (st *ReplicaReadWrite) delete() string {
	id := st.getRandomAsset()
	st.Assets = removeFromSlice(st.Assets, id)

	return fmt.Sprintf("DELETE FROM %s WHERE id='%s';", st.TableName, id)
}

func (st *ReplicaReadWrite) getRandomRegion() string {
	return st.Regions[rand.Intn(len(st.Regions))]
}

func (st *ReplicaReadWrite) getRandomAsset() string {
	return st.Assets[rand.Intn(len(st.Assets))]
}

// generateRow returns a row
func generateRow(rand internal.Random) Row {
	s := uuid.NewString()
	d := randomDate()
	r32 := randomFloat32()
	r64 := randomFloat64()
	return Row{
		ID:                 s,
		Driver_id:          s,
		Connected:          true,
		Created_at:         time.Now(),
		Loc:                [2]float64{r64, r64},
		Bearing:            r32,
		Altitude:           r32,
		Queue_zone_id:      s,
		Queue_zone_left_id: s,
		Queue_zone_left_at: d,
		Queue_zone_hit_at:  d,
		Queue_zone_loc:     [2]float64{r64, r64},
		Loc_ts:             d,
		Speed:              r32,
		Alu_serial:         1234,
		Snoozed_at:         d,
	}
}

func randomDate() time.Time {
	min := time.Date(1970, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2070, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	delta := max - min

	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0)
}

func randomFloat32() float32 {
	max := 180
	min := -180
	return float32(rand.Intn(max-min)+min) + rand.Float32()
}

func randomFloat64() float64 {
	return float64(randomFloat32())
}

func removeFromSlice(list []string, s string) []string {
	index := indexOf(s, list)
	ret := make([]string, 0)
	ret = append(ret, list[:index]...)
	return append(ret, list[index+1:]...)
}

func indexOf(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1 //not found.
}
