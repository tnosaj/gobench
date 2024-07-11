package args_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tnosaj/gobench/pkg/args"
)

func Test(t *testing.T) {
	a := assert.New(t)
	x, err := args.ParseReadWriteSplit("90:10")
	a.NoError(err)
	a.Equal(90, x.Reads)
	a.Equal(10, x.Writes)

	_, err = args.ParseReadWriteSplit("x:10")
	a.EqualError(err, "could not split and convert reads from: x:10")

	_, err = args.ParseReadWriteSplit("90:x")
	a.EqualError(err, "could not split and convert writes from: 90:x")

	_, err = args.ParseReadWriteSplit("DU ARSCH")
	a.EqualError(err, "invalid input \"DU ARSCH\", valid inputs are XX:YY")

}
