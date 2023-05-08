package gpsgen

import (
	"hash/fnv"
	"runtime"

	"github.com/google/uuid"
)

type Option func(*runner)

type Generator struct {
	runners []*runner
}

func New(opts ...Option) *Generator {
	numCpu := runtime.NumCPU()
	gen := Generator{
		runners: make([]*runner, numCpu),
	}
	for i := 0; i < numCpu; i++ {
		gen.runners[i] = newRunner(opts...)
	}
	return &gen
}

func (g *Generator) Attach(dev *Device) {
	g.selectRunner(dev.ID()).attach(dev)
}

func (g *Generator) Detach(id uuid.UUID) {
	g.selectRunner(id).detach(id)
}

func (g *Generator) Run() {
	for i := 0; i < len(g.runners); i++ {
		g.runners[i].run()
	}
}

func (g *Generator) Close() {
	for i := 0; i < len(g.runners); i++ {
		g.runners[i].close()
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
