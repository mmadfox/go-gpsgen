package types

import (
	"github.com/valyala/fastrand"
)

// Random represents a random number generator.
type Random struct {
	min, max int
	rnd      fastrand.RNG
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
	return int(r.rnd.Uint32n(uint32(r.max-r.min))) + r.min
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

	if min == 0 && max == 0 {
		max = 1
	}

	return &Random{
		min: min,
		max: max,
		rnd: fastrand.RNG{},
	}
}
