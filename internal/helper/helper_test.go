package helper_test

import (
	"testing"

	"github.com/tj/assert"
	"gitlab.otters.xyz/jason.tevnan/gobench/internal/db"
	"gitlab.otters.xyz/jason.tevnan/gobench/internal/helper"
)

func TestRows(t *testing.T) {
	a := assert.New(t)
	r := helper.GenerateRow(helper.NewFakeRandomizer())
	a.Equal(db.Row{
		K:   2147483646,
		C:   "ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ",
		Pad: "ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ",
	},
		r)
}
