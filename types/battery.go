package types

import (
	"errors"
	"fmt"

	"github.com/mmadfox/go-gpsgen/curve"
)

const (
	minBatteryVal = 0
	maxBatteryVal = 100
)

var (
	ErrMinBattery = errors.New("types/battery: negative battery value")
	ErrMaxBattery = errors.New("types/battery: invalid battery value")
)

type Battery struct {
	min float64
	max float64
	val float64
	gen *curve.Curve
}

func NewBattery(min, max float64) (*Battery, error) {
	if min < minBatteryVal {
		return nil, ErrMinBattery
	}
	if max > maxBatteryVal || min > max || max <= 0 {
		return nil, ErrMaxBattery
	}
	if min == max {
		max += 10
	}
	if max > maxBatteryVal {
		max = maxBatteryVal
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

func (t *Battery) Min() float64 {
	return t.min
}

func (t *Battery) Max() float64 {
	return t.max
}

func (t *Battery) Value() float64 {
	return t.val
}

func (t *Battery) String() string {
	return fmt.Sprintf("battery: %.2f%%", t.val)
}

func (t *Battery) Next(tick float64) {
	point := t.gen.Point(tick)
	t.val = point.Y
}
