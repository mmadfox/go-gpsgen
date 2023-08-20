package types

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/mmadfox/go-gpsgen/curve"
	"github.com/mmadfox/go-gpsgen/proto"
)

// ErrEmptySensorName indicates that the sensor name is empty.
var ErrEmptySensorName = errors.New("types/sensor: empty name")

// Sensor structure provides a flexible and extensible way to represent and work with sensors.
// It allows for generating different values for different tasks, making it suitable
// for various applications and use cases involving sensor data collection, simulation, or analysis.
type Sensor struct {
	id   string
	name string
	min  float64
	max  float64
	valX float64
	valY float64
	gen  *curve.Curve
}

type (
	SensorSnapshot = proto.Snapshot_SensorType
	SensorMode     = curve.CurveMode
)

const (
	WithSensorRandomMode = curve.ModeDefault
	WithSensorStartMode  = curve.ModeMinStart
	WithSensorEndMode    = curve.ModeMinEnd
)

// NewSensor creates a new Sensor instance with the given name,
// minimum and maximum values, and amplitude.
// The amplitude is used to generate a random curve for the sensor.
//
// Valid amplitude values from 4 to 512.
func NewSensor(name string, min, max float64, amplitude int, mode SensorMode) (*Sensor, error) {
	if err := validateAmplitude(amplitude); err != nil {
		return nil, err
	}
	if len(name) == 0 {
		return nil, ErrEmptySensorName
	}
	if mode == 0 {
		mode = curve.ModeDefault
	}
	gen := curve.New(
		min,
		max,
		amplitude,
		mode,
	)
	return &Sensor{
		id:   uuid.NewString(),
		name: name,
		min:  min,
		max:  max,
		gen:  gen,
	}, nil
}

// ID returns the unique identifier of the sensor.
func (t *Sensor) ID() string {
	return t.id
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

// Shuffle shuffles the generator of the Sensor instance.
func (t *Sensor) Shuffle() {
	t.gen.Shuffle()
}

// Next updates the sensor's values based on the tick value,
// generating new values using the associated generator object.
//
// Tick in the range from 0 to 1.
func (t *Sensor) Next(tick float64) {
	p := t.gen.Point(tick)
	t.valX = p.X
	t.valY = p.Y
}

// Snapshot returns a protobuf snapshot of the Sensor instance.
func (t *Sensor) Snapshot() *SensorSnapshot {
	return &SensorSnapshot{
		Id:   t.id,
		Name: t.name,
		Min:  t.min,
		Max:  t.max,
		ValX: t.valX,
		ValY: t.valY,
		Gen:  t.gen.Snapshot(),
	}
}

// FromSnapshot restores the Sensor instance from a protobuf snapshot.
func (t *Sensor) FromSnapshot(snap *SensorSnapshot) {
	t.id = snap.Id
	t.name = snap.Name
	t.min = snap.Min
	t.max = snap.Max
	t.valX = snap.ValX
	t.valY = snap.ValY
	t.gen = new(curve.Curve)
	t.gen.FromSnapshot(snap.Gen)
}
