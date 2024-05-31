package helper_test

import (
	"testing"

	"github.com/tj/assert"
	"github.com/tnosaj/gobench/internal/helper"
)

func TestFakeRandom(t *testing.T) {
	r := helper.NewFakeRandomizer()
	i := r.Intn(10)
	assert.LessOrEqual(t, i, 10)
}

func TestRealRandom(t *testing.T) {
	r := helper.NewRandomizer()
	i := r.Intn(11)
	assert.LessOrEqual(t, i, 11)
}
