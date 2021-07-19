package helper

import (
	"math/rand"
)

// Randomizer is a real randomizer
type Randomizer struct{}

// NewRandomizer returns a real randomizer
func NewRandomizer() Randomizer {
	return Randomizer{}
}

// Intn is a real Uint
func (r Randomizer) Intn(n int) int {
	return int(rand.Intn(n))
}

// FakeRandomizer is a fake randomizer
type FakeRandomizer struct{}

// NewFakeRandomizer returns a fake randomizer
func NewFakeRandomizer() FakeRandomizer {
	return FakeRandomizer{}
}

// Intn is a fake Uint
func (r FakeRandomizer) Intn(n int) int {
	if n > 10 {
		return n - 1
	}
	return int(rand.Intn(n))
}
