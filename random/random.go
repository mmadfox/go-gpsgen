package random

import (
	"math/rand"
	"sync"
	"time"
)

var defaultRnd = NewRandom()

type Random struct {
	mu  sync.RWMutex
	rnd *rand.Rand
}

func NewRandom() *Random {
	return &Random{
		rnd: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (r *Random) Intn(n int) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.rnd.Intn(n)
}

func (r *Random) Float64() float64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.rnd.Float64()
}

func (r *Random) ExpFloat64() float64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.rnd.ExpFloat64()
}
