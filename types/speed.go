package types

import (
	"errors"
	"fmt"

	"github.com/mmadfox/go-gpsgen/curve"
	"github.com/mmadfox/go-gpsgen/proto"
)

const (
	MaxSpeedVal = 1000
	MinSpeedVal = 0
)

var (
	// ErrMinSpeed indicates that the speed value is less than 0.
	ErrMinSpeed = errors.New("types/speed: value is less than 0")
	// ErrMaxSpeed indicates that the speed value is greater than 1000.
	ErrMaxSpeed = errors.New("types/speed: value of greater than 1000")
	// ErrSpeedMinGreaterMax indicates that the minimum speed value is greater than the maximum speed value.
	ErrSpeedMinGreaterMax = errors.New("types/speed: min value greater than max value")
)

// Speed represents a speed value with a minimum and maximum range and contains a value generator.
type Speed struct {
	min float64
	max float64
	val float64
	gen *curve.Curve
}

// NewSpeed creates a new Speed instance with the specified minimum and maximum values.
// It also takes an amplitude parameter for generating a random curve.
// The minimum value is 0 to maximum 1000, and the amplitude parameter must be 4 to 512.
func NewSpeed(min, max float64, amplitude int) (*Speed, error) {
	if min < MinSpeedVal {
		return nil, ErrMinSpeed
	}
	if max > MaxSpeedVal {
		return nil, ErrMaxSpeed
	}
	if max <= 0 {
		max = min
	}
	if min > max {
		return nil, ErrSpeedMinGreaterMax
	}
	if min == 0 && max == 0 {
		min = 1
		max = 1
	}
	if err := validateAmplitude(amplitude); err != nil {
		return nil, err
	}
	gen := curve.New(
		min,
		max,
		amplitude,
		curve.ModeDefault,
	)
	return &Speed{
		min: min,
		max: max,
		gen: gen,
	}, nil
}

// Min returns the minimum speed value of the Speed instance.
func (t *Speed) Min() float64 {
	return t.min
}

// Max returns the maximum speed value of the Speed instance.
func (t *Speed) Max() float64 {
	return t.max
}

// Value returns the current speed value of the Speed instance.
func (t *Speed) Value() float64 {
	if t.isStatic() {
		return t.min
	}
	return t.val
}

// String returns a formatted string representation of the speed value.
func (t *Speed) String() string {
	return fmt.Sprintf("speed: %.2f m/s", t.val)
}

// Shuffle shuffles the generator of the Speed instance.
func (t *Speed) Shuffle() {
	t.gen.Shuffle()
}

// Next generates the next speed value based on a given tick value,
// using the underlying random curve generator.
//
// Tick in the range from 0 to 1.
func (t *Speed) Next(tick float64) {
	if t.isStatic() {
		return
	}
	p := t.gen.Point(tick)
	t.val = p.Y
}

// Snapshot returns a protobuf snapshot of the Speed instance.
func (t *Speed) Snapshot() *proto.Snapshot_CommonType {
	return &proto.Snapshot_CommonType{
		Min: t.min,
		Max: t.max,
		Val: t.val,
		Gen: t.gen.Snapshot(),
	}
}

// FromSnapshot restores the Speed instance from a protobuf snapshot.
func (t *Speed) FromSnapshot(snap *proto.Snapshot_CommonType) {
	t.min = snap.Min
	t.max = snap.Max
	t.val = snap.Val
	t.gen = new(curve.Curve)
	t.gen.FromSnapshot(snap.Gen)
}

func (t *Speed) isStatic() bool {
	return t.min == t.max
}
