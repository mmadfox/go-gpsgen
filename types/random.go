package types

import (
	"math/rand"
	"time"
)

type Random struct {
	min, max int
	rnd      *rand.Rand
}

func (r *Random) Value() int {
	return r.rnd.Intn(r.max-r.min) + r.min
}

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
