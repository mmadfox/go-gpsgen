package gpsgen

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type runner struct {
	mu        sync.RWMutex
	devices   map[uuid.UUID]*Device
	ticker    *time.Ticker
	closeCh   chan struct{}
	closeOnce sync.Once
	interval  time.Duration
}

func NextTick(d ...*Device) {
	if len(d) == 0 {
		return
	}
	if len(d) == 1 {
		d[0].nextTick(1)
		return
	}
	for i := 0; i < len(d); i++ {
		d[i].nextTick(1)
	}
}

func WithInterval(dur time.Duration) Option {
	return func(r *runner) {
		r.interval = dur
	}
}

func newRunner(opts ...Option) *runner {
	r := &runner{
		devices:  make(map[uuid.UUID]*Device, 64),
		closeCh:  make(chan struct{}),
		interval: time.Second,
	}
	for i := 0; i < len(opts); i++ {
		opts[i](r)
	}
	return r
}

func (r *runner) attach(d *Device) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.devices[d.ID()]; ok {
		return
	}
	go d.handleChange()
	r.devices[d.ID()] = d
}

func (r *runner) detach(id uuid.UUID) {
	r.mu.Lock()
	defer r.mu.Unlock()
	dev, ok := r.devices[id]
	if !ok {
		return
	}
	close(dev.stateCh)
	delete(r.devices, id)
}

func (r *runner) run() {
	r.ticker = time.NewTicker(r.interval)
	var tick float64
	if r.interval.Seconds() < 1 {
		tick = 1
	} else {
		tick = r.interval.Seconds()
	}
	go func() {
		for {
			select {
			case <-r.closeCh:
				r.ticker.Stop()
				return
			case <-r.ticker.C:
				r.mu.RLock()
				for _, dev := range r.devices {
					dev.nextTick(tick)
				}
				r.mu.RUnlock()
			}
		}
	}()
}

func (r *runner) close() {
	r.closeOnce.Do(func() {
		close(r.closeCh)
	})
}
