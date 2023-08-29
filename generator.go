package gpsgen

import (
	"context"
	"math"
	"runtime"
	"sync"
	"time"

	pb "github.com/mmadfox/go-gpsgen/proto"
	"google.golang.org/protobuf/proto"
)

// Generator represents the GPS data generator.
type Generator struct {
	mu         sync.RWMutex
	dmu        sync.RWMutex
	devices    []*Device
	wg         sync.WaitGroup
	packet     *pb.Packet
	cancelFunc context.CancelFunc
	ctx        context.Context
	sliceCh    chan slice
	index      int
	ticker     *time.Ticker
	waitCh     chan struct{}
	nextCh     chan struct{}
	numWorkers int
	packetSize int
	workerSize int
	interval   time.Duration
	onError    func(error)
	onPacket   func([]byte)
	onNext     func()
}

// Options defines the configuration options for the Generator.
type Options struct {
	// Interval determines the time interval between data generation iterations. Default three seconds.
	Interval time.Duration

	// PacketSize specifies the size of the data packet generated per iteration. Default 8192.
	PacketSize int

	// NumWorkers sets the number of concurrent workers for data processing. Default 4.
	NumWorkers int
}

// NewOptions creates a new Options instance with default values.
func NewOptions() *Options {
	return &Options{
		Interval:   3*time.Second,
		PacketSize: 8192,
		NumWorkers: 4,
	}
}

func (o *Options) prepare() {
	if o.PacketSize < 1024 {
		o.PacketSize = 1024
	}
	if o.Interval <= 0 {
		o.Interval = time.Second
	}
	if o.NumWorkers <= 0 {
		o.NumWorkers = runtime.NumCPU()
	}
}

// New creates a new GPS data generator with the provided options.
func New(opts *Options) *Generator {
	if opts == nil {
		opts = NewOptions()
	}

	opts.prepare()

	gen := Generator{
		interval:   opts.Interval,
		packetSize: opts.PacketSize,
		numWorkers: opts.NumWorkers,
	}
	gen.workerSize = gen.packetSize / gen.numWorkers
	gen.packet = &pb.Packet{Devices: make([]*pb.Device, gen.packetSize)}
	gen.devices = make([]*Device, 0, gen.packetSize)
	gen.sliceCh = make(chan slice, gen.numWorkers)
	gen.waitCh = make(chan struct{}, gen.numWorkers)
	gen.nextCh = make(chan struct{}, 1)
	gen.ticker = time.NewTicker(gen.interval)
	gen.ctx, gen.cancelFunc = context.WithCancel(context.Background())
	return &gen
}

// HasTracker checks if a tracker with the given deviceID exists in the Generator.
func (g *Generator) HasTracker(deviceID string) bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	_, ok := g.find(deviceID)
	return ok
}

// Attach attaches the provided device to the generator.
func (g *Generator) Attach(d *Device) error {
	if d == nil {
		return nil
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	if err := d.mount(); err != nil {
		return err
	}

	g.devices = append(g.devices, d)
	return nil
}

// Detach detaches a device with the given ID from the generator.
func (g *Generator) Detach(deviceID string) error {
	if len(deviceID) == 0 {
		return nil
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	return g.delete(deviceID)
}

// Each iterates over the collection of Device objects managed by the Generator
// and applies the provided function to each device.
func (g *Generator) Each(fn func(int, *Device) bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	for i := 0; i < len(g.devices); i++ {
		if next := fn(i, g.devices[i]); !next {
			break
		}
	}
}

// Lookup searches for a device with the given ID
// and returns it along with a boolean indicating its existence.
func (g *Generator) Lookup(deviceID string) (*Device, bool) {
	if len(deviceID) == 0 {
		return nil, false
	}

	g.mu.RLock()
	defer g.mu.RUnlock()

	return g.find(deviceID)
}

// NumDevices returns the number of devices attached to the generator.
func (g *Generator) NumDevices() int {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return len(g.devices)
}

// OnError sets a callback function to handle errors during data generation.
func (g *Generator) OnError(fn func(error)) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.onError = fn
}

// OnPacket sets a callback function to handle generated data packets.
func (g *Generator) OnPacket(fn func([]byte)) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.onPacket = fn
}

// OnNext sets a callback function to be executed at each "next step" of data generation.
func (g *Generator) OnNext(fn func()) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.onNext = fn
}

// Close stops the data generation process and closes the generator.
func (g *Generator) Close() {
	g.cancelFunc()
}

// Run starts the data generation process using configured settings and attached devices.
func (g *Generator) Run() {
	g.wg.Add(g.numWorkers)
	for i := 0; i < g.numWorkers; i++ {
		go g.doWorker(i + 1)
	}

	g.wg.Add(1)
	go g.run()

	g.wg.Add(1)
	go g.doNext()

	g.wg.Wait()
}

func (g *Generator) run() {
	defer func() {
		g.wg.Done()
		g.ticker.Stop()
	}()

	lastTime := time.Now()
	var tick float64

loop:
	for {
		select {
		case <-g.ctx.Done():
			break loop
		case t := <-g.ticker.C:
			if g.isClosed() {
				break loop
			}

			tick = math.Abs(lastTime.Sub(t).Seconds())
			if tick < 1 {
				tick = 1
			}

			g.mu.RLock()
			for i := 0; i < len(g.devices); i++ {
				g.devices[i].Next(tick)
			}

			g.index = 0
			for i := 0; i < len(g.devices); i++ {
				g.dmu.Lock()
				g.packet.Devices[g.index] = g.devices[i].State()
				g.dmu.Unlock()

				if g.index+1 == g.packetSize {
					g.mu.RUnlock()
					g.flush()
					g.mu.RLock()
					g.index = 0
				} else {
					g.index++
				}
			}
			g.mu.RUnlock()

			g.flush()
			g.resetPacket()
			g.notifyNextTick()

			lastTime = t
		}
	}
}

type slice struct {
	from, to int
}

func (g *Generator) doWorker(n int) {
	defer func() {
		g.wg.Done()
		g.waitCh <- struct{}{}
	}()
	
	pck := &pb.Packet{}
	enc := proto.MarshalOptions{}

	for {
		select {
		case <-g.ctx.Done():
			return
		case s := <-g.sliceCh:
			if g.isClosed() {
				return
			}

			if s.from == 0 && s.to == 0 {
				g.waitCh <- struct{}{}
				continue
			}

			pck.Timestamp = time.Now().Unix()

			buf := make([]byte, 0, g.workerSize)
			g.dmu.RLock()
			pck.Devices = g.packet.Devices[s.from:s.to]
			data, err := enc.MarshalAppend(buf, pck)
			g.dmu.RUnlock()

			if err != nil && g.onError != nil {
				g.onError(err)
			} else if g.onPacket != nil {
				g.onPacket(data)
			}

			if g.isClosed() {
				return
			}

			g.waitCh <- struct{}{}
		}
	}
}

func (g *Generator) notifyNextTick() {
	select {
	case g.nextCh <- struct{}{}:
	default:
	}
}

func (g *Generator) doNext() {
	defer g.wg.Done()

	for {
		select {
		case <-g.ctx.Done():
			return
		case <-g.nextCh:
			if g.isClosed() {
				return
			}
			if g.onNext == nil {
				continue
			}
			g.onNext()
		}
	}
}

func (g *Generator) flush() {
	if g.index <= g.numWorkers+g.workerSize {
		g.sliceCh <- slice{from: 0, to: g.index}
		if g.isClosed() {
			return
		}
		<-g.waitCh
	} else {
		step := g.index / g.numWorkers
		rem := g.index % g.numWorkers
		for i := 0; i < g.numWorkers; i++ {
			if g.isClosed() {
				return
			}
			var s slice
			s.from = i * step
			s.to = i*step + step
			if i == g.numWorkers-1 {
				s.to += rem
			}
			g.sliceCh <- s
		}
		for i := 0; i < g.numWorkers; i++ {
			if g.isClosed() {
				return
			}
			<-g.waitCh
		}
	}
}

func (g *Generator) isClosed() bool {
	select {
	case <-g.ctx.Done():
		return true
	default:
		return false
	}
}

func (g *Generator) resetPacket() {
	g.dmu.Lock()
	defer g.dmu.Unlock()
	for i := 0; i < len(g.packet.Devices); i++ {
		g.packet.Devices[i] = nil
	}
}

func (g *Generator) delete(deviceID string) error {
	for i := 0; i < len(g.devices); i++ {
		if g.devices[i].ID() == deviceID {
			err := g.devices[i].unmount()
			g.devices = append(g.devices[:i], g.devices[i+1:]...)
			return err
		}
	}
	return nil
}

func (g *Generator) find(deviceID string) (d *Device, ok bool) {
	for i := 0; i < len(g.devices); i++ {
		if g.devices[i].ID() == deviceID {
			d = g.devices[i]
			ok = true
			break
		}
	}
	return
}
