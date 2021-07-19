package args

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type ReadWriteSplit struct {
	Reads  int
	Writes int
}

// ParseReadWriteSplit gets read:write
func ParseReadWriteSplit(split string) (ReadWriteSplit, error) {
	s := strings.Split(split, ":")
	if len(s) != 2 {
		return ReadWriteSplit{}, fmt.Errorf("invalid input %q, valid inputs are XX:YY", split)
	}

	r, err := strconv.ParseFloat(strings.Split(split, ":")[0], 64)
	if err != nil {
		return ReadWriteSplit{}, fmt.Errorf("could not split and convert reads from: %s", split)
	}

	w, err := strconv.ParseFloat(strings.Split(split, ":")[1], 64)
	if err != nil {
		return ReadWriteSplit{}, fmt.Errorf("could not split and convert writes from: %s", split)
	}

	// make a percentage
	rp := int(math.Round(100 * (r / (r + w))))
	wp := int(math.Round(100 * (w / (r + w))))

	return ReadWriteSplit{Reads: rp, Writes: wp}, nil
}
