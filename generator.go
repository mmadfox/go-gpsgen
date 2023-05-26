package gpsgen

import (
	"hash/fnv"
	"runtime"
	"sync"

	"github.com/google/uuid"
)

// Option uses to configure the behavior of the underlying
// runner instances when creating a new Generator instance.
type Option func(*runner)

// Generator provides a way to manage and control multiple runner instances,
// allowing for concurrent execution and management of attached devices.
type Generator struct {
	OnClose func()

	runners []*runner
	once    sync.Once
}

// New creates a new Generator instance.
//
// It accepts Option arguments to configure the underlying runner instances.
// It creates a slice of runner instances with a length equal to the number of available CPUs.
func New(opts ...Option) *Generator {
	numCpu := runtime.NumCPU()
	if numCpu > 4 {
		numCpu = 4
	}
	gen := Generator{
		runners: make([]*runner, numCpu),
	}
	for i := 0; i < numCpu; i++ {
		gen.runners[i] = newRunner(opts...)
	}
	return &gen
}

// Attach attaches a Device to the Generator by selecting the appropriate runner
// based on the device's ID and calling the attach method of the selected runner.
func (g *Generator) Attach(dev *Device) {
	g.selectRunner(dev.ID()).attach(dev)
}

// Lookup retrieves a Device object from the Generator based on its deviceID,
// leveraging the selected runner responsible for handling the device.
func (g *Generator) Lookup(deviceID uuid.UUID) (*Device, error) {
	return g.selectRunner(deviceID).lookup(deviceID)
}

// Detach detaches a Device from the Generator by selecting the appropriate runner
// based on the device's ID and calling the detach method of the selected runner.
func (g *Generator) Detach(id uuid.UUID) {
	g.selectRunner(id).detach(id)
}

// Run runs the run method of each runner in the Generator.
func (g *Generator) Run() {
	for i := 0; i < len(g.runners); i++ {
		g.runners[i].run()
	}
}

// Close calls the close method of each runner in the Generator.
func (g *Generator) Close() {
	for i := 0; i < len(g.runners); i++ {
		g.runners[i].close()
	}
	if g.OnClose != nil {
		g.once.Do(func() {
			g.OnClose()
		})
	}
}

func (g *Generator) selectRunner(deviceID uuid.UUID) *runner {
	return g.runners[g.runnerIndex(deviceID)]
}

func (g *Generator) runnerIndex(deviceID uuid.UUID) uint32 {
	a, _ := deviceID.MarshalBinary()
	h := fnv.New32a()
	_, _ = h.Write(a)
	return h.Sum32() % uint32(len(g.runners))
}
