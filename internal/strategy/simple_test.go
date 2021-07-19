package strategy_test

import (
	"testing"

	"github.com/tj/assert"
	"gitlab.otters.xyz/jason.tevnan/gobench/internal"
	"gitlab.otters.xyz/jason.tevnan/gobench/internal/helper"
	"gitlab.otters.xyz/jason.tevnan/gobench/internal/strategy"
	"gitlab.otters.xyz/jason.tevnan/gobench/pkg/args"
)

func TestReads(t *testing.T) {
	a := assert.New(t)

	r := strategy.SimpleReadWrite{
		S: internal.Settings{
			TableName: "testtest",
			ReadWriteSplit: args.ReadWriteSplit{
				Reads:  100,
				Writes: 0,
			},
			Randomizer: helper.NewFakeRandomizer(),
		},
		MaxIDCount: 1,
	}
	for i := 0; i <= 10; i++ {
		r, l := r.CreateCommand()
		switch l {
		case "getPk":
			a.Equal("select id,k,c,pad from testtest where id=0;", r)
		case "getSk":
			a.Equal("select id,k,c,pad from testtest where k=2147483646;", r)
		default:
			a.Fail("Unsupported Label received" + l)
		}
	}
}
func TestWrites(t *testing.T) {
	b := assert.New(t)

	w := strategy.SimpleReadWrite{
		S: internal.Settings{
			TableName: "testtest",
			ReadWriteSplit: args.ReadWriteSplit{
				Reads:  -1,
				Writes: 100,
			},
			Randomizer: helper.NewFakeRandomizer(),
		},
		MaxIDCount: 1,
	}
	for i := 0; i <= 10; i++ {
		r, l := w.CreateCommand()
		switch l {
		case "create":
			b.Equal("INSERT INTO testtest(k, c , pad) VALUES (2147483646,'ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ','ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ');", r)
		case "delete":
			b.Equal("DELETE FROM testtest WHERE id=0;", r)
		case "update":
			b.Equal("UPDATE testtest SET c='ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ', pad='ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ' WHERE id=0;", r)
		default:
			b.Fail("Unsupported Label received: " + l)
		}
	}

}
