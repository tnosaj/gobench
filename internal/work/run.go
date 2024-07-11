package work

import (
	"log"

	"github.com/sirupsen/logrus"
	"github.com/tnosaj/gobench/internal"
	"github.com/tnosaj/gobench/internal/strategy"
)

func run(s *internal.Settings, wp *workerPool, st strategy.ExecutionStrategy) {
	logrus.Infof("run")

	logrus.Infof("Running with a %d:%d::read:write split and strategy: %s", s.ReadWriteSplit.Reads, s.ReadWriteSplit.Writes, s.Strategy)

	// Catch other strategies
	s.ServerStatus = "busy"
	var runner ExecutionType
	switch s.DurationType {
	case "events":
		runner = RunForEventCount{
			s:  s,
			wp: wp,
			st: st,
		}
	case "seconds":
		log.Fatalf("Sorry, seconds is not implemented yet")
	}

	runner.Run()
	logrus.Infof("Done")
	s.ServerStatus = "free"
}

// RunForEventCount do stuffs
type RunForEventCount struct {
	s  *internal.Settings
	wp *workerPool
	st strategy.ExecutionStrategy
}

// Run for a number of events
func (r RunForEventCount) Run() {

	for i := 0; i < r.s.Duration; i++ {
		r.wp.do(func() {
			r.st.RunCommand()
		})
	}
	r.wp.stop()

}
