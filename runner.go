package gpsgen

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type runner struct {
	mu        sync.RWMutex
	devices   []*Device
	index     map[uuid.UUID]*Device
	buffer    []*Device
	ticker    *time.Ticker
	closeCh   chan struct{}
	closeOnce sync.Once
	interval  time.Duration
	running   bool
}

func NextTick(d ...*Device) {
	if len(d) == 0 {
		return
	}
	if len(d) == 1 {
		d[0].nextTick(1, 1)
		return
	}
	for i := 0; i < len(d); i++ {
		d[i].nextTick(1, 1)
	}
}

func WithInterval(dur time.Duration) Option {
	return func(r *runner) {
		r.interval = dur
	}
}

func newRunner(opts ...Option) *runner {
	r := &runner{
		closeCh:  make(chan struct{}),
		interval: time.Second,
		devices:  make([]*Device, 0, 64),
		index:    make(map[uuid.UUID]*Device, 64),
		buffer:   make([]*Device, 0, 8),
	}
	for i := 0; i < len(opts); i++ {
		opts[i](r)
	}
	return r
}

func (r *runner) attach(d *Device) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.running {
		r.buffer = append(r.buffer, d)
	} else {
		_, ok := r.index[d.id]
		if ok {
			return
		}
		r.index[d.id] = d
		r.devices = append(r.devices, d)
		d.mount()
	}
}

func (r *runner) lookup(id uuid.UUID) (*Device, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	dev, ok := r.index[id]
	if !ok {
		return nil, fmt.Errorf("gpsgen: device %s not found", id)
	}
	return dev, nil
}

func (r *runner) detach(id uuid.UUID) {
	r.mu.Lock()
	defer r.mu.Unlock()

	dev, ok := r.index[id]
	if !ok {
		return
	}

	for i := 0; i < len(r.devices); i++ {
		if r.devices[i].id == id {
			r.devices = append(r.devices[:i], r.devices[i+1:]...)
			break
		}
	}

	if len(r.buffer) > 0 {
		for i := 0; i < len(r.buffer); i++ {
			if r.buffer[i] == nil {
				continue
			}
			if r.buffer[i].id == id {
				r.buffer = append(r.buffer[:i], r.buffer[i+1:]...)
				break
			}
		}
	}

	delete(r.index, id)
	dev.unmount()
}

func (r *runner) run() {
	r.ticker = time.NewTicker(r.interval)
	r.running = true

	go func() {
		defer func() {
			r.closeDevices()
		}()

		var prevTick time.Time
	loop:
		for {
			select {
			case <-r.closeCh:
				r.ticker.Stop()
				return
			case ts := <-r.ticker.C:
				select {
				case <-r.closeCh:
					break loop
				default:
				}

				tick := ts.Sub(prevTick).Seconds()
				if prevTick.Second() > 0 && tick > 5 {
					fmt.Printf("Warning! processing lag %f seconds instead of 1 second\n", tick)
				}

				for _, dev := range r.devices {
					dev.nextTick(1, tick)
				}
				prevTick = ts

				r.flushBuffer()
			}
		}
	}()
}

func (r *runner) closeDevices() {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, dev := range r.index {
		dev.unmount()
	}

	r.devices = nil
	r.index = nil
	r.buffer = nil
}

func (r *runner) flushBuffer() {
	if len(r.buffer) == 0 {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := 0; i < len(r.buffer); i++ {
		dev := r.buffer[i]
		_, ok := r.index[dev.id]
		if ok {
			continue
		}
		r.index[dev.id] = dev
		r.devices = append(r.devices, dev)
		r.buffer[i] = nil
		dev.mount()
	}
	if cap(r.buffer) > 1024 {
		r.buffer = make([]*Device, 0, 8)
	} else {
		r.buffer = r.buffer[:0]
	}
}

func (r *runner) close() {
	r.closeOnce.Do(func() {
		close(r.closeCh)
	})
}
