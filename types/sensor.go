package types

import (
	"errors"
	"fmt"

	"github.com/mmadfox/go-gpsgen/curve"
	"github.com/mmadfox/go-gpsgen/proto"
)

// ErrEmptySensorName indicates that the sensor name is empty.
var ErrEmptySensorName = errors.New("types/sensor: empty name")

// Sensor structure provides a flexible and extensible way to represent and work with sensors.
// It allows for generating different values for different tasks, making it suitable
// for various applications and use cases involving sensor data collection, simulation, or analysis.
type Sensor struct {
	name string
	min  float64
	max  float64
	valX float64
	valY float64
	gen  *curve.Curve
}

// NewSensor creates a new Sensor instance with the given name,
// minimum and maximum values, and amplitude.
// The amplitude is used to generate a random curve for the sensor.
//
// Valid amplitude values from 4 to 512.
func NewSensor(name string, min, max float64, amplitude int) (*Sensor, error) {
	if err := validateAmplitude(amplitude); err != nil {
		return nil, err
	}
	if len(name) == 0 {
		return nil, ErrEmptySensorName
	}
	gen, err := curve.RandomCurveWithMode(
		min,
		max,
		amplitude,
		curve.ModeDefault,
	)
	if err != nil {
		return nil, err
	}
	return &Sensor{
		name: name,
		min:  min,
		max:  max,
		gen:  gen,
	}, nil
}

// ToProto converts the sensor object to its corresponding
// protobuf message representation.
func (t *Sensor) ToProto() *proto.SensorState {
	return &proto.SensorState{
		Min:  t.min,
		Max:  t.max,
		ValX: t.valX,
		ValY: t.valY,
		Name: t.name,
		Gen:  t.gen.ToProto(),
	}
}

// FromProto sets the fields of the sensor object based on the values
// from the provided protobuf message.
func (t *Sensor) FromProto(sensor *proto.SensorState) {
	t.gen = new(curve.Curve)
	t.gen.FromProto(sensor.Gen)
	t.min = sensor.Min
	t.max = sensor.Max
	t.valX = sensor.ValX
	t.valY = sensor.ValY
	t.name = sensor.Name
}

// Name returns the name of the sensor.
func (t *Sensor) Name() string {
	return t.name
}

// Min returns the minimum value of the sensor.
func (t *Sensor) Min() float64 {
	return t.min
}

// Max returns the maximum value of the sensor.
func (t *Sensor) Max() float64 {
	return t.max
}

// ValueX returns the current value of the sensor along the X-axis.
func (t *Sensor) ValueX() float64 {
	return t.valX
}

// ValueY returns the current value of the sensor along the Y-axis.
func (t *Sensor) ValueY() float64 {
	return t.valY
}

// String returns a formatted string representation of the sensor object.
func (t *Sensor) String() string {
	return fmt.Sprintf("%s: valX=%.8f, valY=%.8f", t.name, t.valX, t.valY)
}

// Next updates the sensor's values based on the tick value,
// generating new values using the associated generator object.
//
// Tick in the range from 0 to 1.
func (t *Sensor) Next(tick float64) {
	point := t.gen.Point(tick)
	t.valX = point.X
	t.valY = point.Y
}
