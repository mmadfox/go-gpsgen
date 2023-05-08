package types

import (
	"fmt"

	"github.com/mmadfox/go-gpsgen/curve"
)

type Sensor struct {
	name string
	min  float64
	max  float64
	valX float64
	valY float64
	gen  *curve.Curve
}

func NewSensor(name string, min, max float64, amplitude int) (*Sensor, error) {
	if amplitude > 1024 {
		amplitude = 1024
	}
	gen, err := curve.RandomCurveWithMode(min, max, amplitude, curve.ModeDefault|curve.ModeMinStart|curve.ModeMinEnd)
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

func (t *Sensor) Name() string {
	return t.name
}

func (t *Sensor) Min() float64 {
	return t.min
}

func (t *Sensor) Max() float64 {
	return t.max
}

func (t *Sensor) ValueX() float64 {
	return t.valX
}

func (t *Sensor) ValueY() float64 {
	return t.valY
}

func (t *Sensor) String() string {
	return fmt.Sprintf("%s: valX=%.8f, valY=%.8f", t.name, t.valX, t.valY)
}

func (t *Sensor) Next(tick float64) {
	point := t.gen.Point(tick)
	t.valX = point.X
	t.valY = point.Y
}
