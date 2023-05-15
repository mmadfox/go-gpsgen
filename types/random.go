package types

import (
	"math/rand"
	"time"
)

// Random represents a random number generator.
type Random struct {
	min, max int
	rnd      *rand.Rand
}

// Min returns the minimum value.
func (r *Random) Min() int {
	return r.min
}

// Max returns the maximum value.
func (r *Random) Max() int {
	return r.max
}

// Value generates and returns a random integer value within the range specified by min and max.
func (r *Random) Value() int {
	return r.rnd.Intn(r.max-r.min) + r.min
}

// NewRandom creates a new Random instance with the specified minimum and maximum values.
func NewRandom(min, max int) *Random {
	if min < 0 {
		min = 0
	}
	if max < 0 {
		min = 0
		max = 1
	}
	return &Random{
		min: min,
		max: max,
		rnd: rand.New(rand.NewSource(time.Now().UnixMicro())),
	}
}
