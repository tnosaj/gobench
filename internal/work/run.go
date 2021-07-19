package work

import (
	"log"

	"github.com/sirupsen/logrus"
	"gitlab.otters.xyz/jason.tevnan/gobench/internal"
	"gitlab.otters.xyz/jason.tevnan/gobench/internal/strategy"
)

func run(s internal.Settings, wp *workerPool) {
	logrus.Infof("run")
	maxIDCount := getMaxIDCount(s)
	logrus.Infof("Query from 0 to %d", maxIDCount)
	logrus.Infof("Running with a %d:%d::read:write split", s.ReadWriteSplit.Reads, s.ReadWriteSplit.Writes)

	// Catch other strategies
	var st internal.ExecutionStrategy
	switch s.Strategy {
	case "simple":
		st = strategy.SimpleReadWrite{
			S:          s,
			MaxIDCount: maxIDCount,
		}
	}

	var runner internal.ExecutionType
	switch s.DurationType {
	case "events":
		runner = RunForEventCount{
			s:          s,
			maxIDCount: maxIDCount,
			wp:         wp,
			st:         st,
		}
	case "seconds":
		log.Fatalf("Sorry, seconds is not implemented yet")
	}

	runner.Run()
}

func getMaxIDCount(s internal.Settings) int {

	count, err := s.DBInterface.ExecStatementWithReturnInt("select count(id) from " + s.TableName + ";")

	if err != nil {
		logrus.Fatalf("could not get max id count with error: %q", err)
	}
	return count

}

// RunForEventCount do stuffs
type RunForEventCount struct {
	s          internal.Settings
	maxIDCount int
	wp         *workerPool
	st         internal.ExecutionStrategy
}

// Run for a number of events
func (r RunForEventCount) Run() {

	for i := 0; i < r.s.Duration; i++ {
		r.wp.do(func() {
			r.s.DBInterface.ExecStatement(r.st.CreateCommand())
		})
	}
	r.wp.stop()

}
