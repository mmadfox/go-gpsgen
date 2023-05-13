package types

import (
	"errors"
	"fmt"

	"github.com/mmadfox/go-gpsgen/curve"
	"github.com/mmadfox/go-gpsgen/proto"
)

const (
	minBatteryVal = 0
	maxBatteryVal = 100
)

var (
	// ErrMinBattery indicates that the value is less than 0%.
	ErrMinBattery = errors.New("types/battery: value is less than 0%")
	// ErrMaxBattery indicates that the value greater than 100%.
	ErrMaxBattery = errors.New("types/battery: value greater than 100%")
)

// Battery represents a Battery object with various methods and fields.
type Battery struct {
	min float64
	max float64
	val float64
	gen *curve.Curve
}

// NewBattery creates a new battery instance and returns a pointer to it.
// It validates the provided minimum and maximum values against predefined limits
// and returns the corresponding errors if the values are out of range.
// Valid values from 0 to 100%.
func NewBattery(min, max float64) (*Battery, error) {
	if min < minBatteryVal {
		return nil, ErrMinBattery
	}
	if max > maxBatteryVal {
		return nil, ErrMaxBattery
	}
	if min > max {
		min = max
	}
	gen, err := curve.RandomCurveWithMode(min, max, 4, curve.ModeMaxMin)
	if err != nil {
		return nil, err
	}
	return &Battery{
		min: min,
		max: max,
		gen: gen,
	}, nil
}

// ToProto converts the battery object into a protobuf message
// (proto.TypeState) and returns it.
func (t *Battery) ToProto() *proto.TypeState {
	return &proto.TypeState{
		Min: t.min,
		Max: t.max,
		Val: t.val,
		Gen: t.gen.ToProto(),
	}
}

// FromProto sets the battery object's fields based on the values
// from a protobuf message (proto.TypeState).
func (t *Battery) FromProto(battery *proto.TypeState) {
	t.gen = new(curve.Curve)
	t.gen.FromProto(battery.Gen)
	t.min = battery.Min
	t.max = battery.Max
	t.val = battery.Val
}

// Min returns the minimum value of the battery.
func (t *Battery) Min() float64 {
	return t.min
}

// Max returns the maximum value of the battery.
func (t *Battery) Max() float64 {
	return t.max
}

// Value returns the current value of the battery.
func (t *Battery) Value() float64 {
	return t.val
}

// String returns a formatted string representation of the battery object,
// displaying the battery level as a percentage.
func (t *Battery) String() string {
	return fmt.Sprintf("battery: %.2f%%", t.val)
}

// Next updates the battery object's value based on the tick
// value by using the generator object.
//
// Tick in the range from 0 to 1.
func (t *Battery) Next(tick float64) {
	point := t.gen.Point(tick)
	t.val = point.Y
}
