package types

import (
	"errors"
	"fmt"

	"github.com/mmadfox/go-gpsgen/curve"
	"github.com/mmadfox/go-gpsgen/proto"
)

const (
	maxSpeedVal = 1000
	minSpeedVal = 0
)

var (
	// ErrMinSpeed indicates that the speed value is less than 0.
	ErrMinSpeed = errors.New("types/speed: value is less than 0")
	// ErrMaxSpeed indicates that the speed value is greater than 1000.
	ErrMaxSpeed = errors.New("types/speed: value of greater than 1000")
	// ErrSpeedMinGreaterMax indicates that the minimum speed value is greater than the maximum speed value.
	ErrSpeedMinGreaterMax = errors.New("types/speed: min value greater than max value")
)

// Speed represents a speed value with a
// minimum and maximum range and contains value generator.
// Basic unit of speed measurement is meters per second.
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
	if min < minSpeedVal {
		return nil, ErrMinSpeed
	}
	if max > maxSpeedVal {
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
	gen, err := curve.RandomCurveWithMode(
		min,
		max,
		amplitude,
		curve.ModeDefault|curve.ModeMinStart|curve.ModeMinEnd,
	)
	if err != nil {
		return nil, err
	}
	return &Speed{
		min: min,
		max: max,
		gen: gen,
	}, nil
}

// ToProto converts the Speed instance into a
// proto representation (proto.TypeState) and returns a pointer to it.
func (t *Speed) ToProto() *proto.TypeState {
	return &proto.TypeState{
		Min: t.min,
		Max: t.max,
		Val: t.val,
		Gen: t.gen.ToProto(),
	}
}

// FromProto sets the values of the Speed instance based
// on the provided proto representation (proto.TypeState).
func (t *Speed) FromProto(speed *proto.TypeState) {
	t.gen = new(curve.Curve)
	t.gen.FromProto(speed.Gen)
	t.min = speed.Min
	t.max = speed.Max
	t.val = speed.Val
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
	return t.val
}

// String returns a formatted string representation of the speed value.
func (t *Speed) String() string {
	return fmt.Sprintf("speed: %.2f m/s", t.val)
}

// Next generates the next speed value based on a given tick value,
// using the underlying random curve generator.
//
// Tick in the range from 0 to 1.
func (t *Speed) Next(tick float64) {
	point := t.gen.Point(tick)
	t.val = point.Y
}
