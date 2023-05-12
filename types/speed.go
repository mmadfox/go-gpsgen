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
	ErrMinSpeed = errors.New("types/speed: negative speed value")
	ErrMaxSpeed = errors.New("types/speed: invalid speed value")
)

type Speed struct {
	min float64
	max float64
	val float64
	gen *curve.Curve
}

func NewSpeed(min, max float64, amplitude int) (*Speed, error) {
	if min < minSpeedVal {
		return nil, ErrMinSpeed
	}
	if max > maxSpeedVal || min > max || max <= 0 {
		return nil, ErrMaxSpeed
	}
	if min == max {
		max += 3
	}
	if amplitude < 4 {
		amplitude = 4
	}
	if amplitude > 1024 {
		amplitude = 1024
	}
	gen, err := curve.RandomCurveWithMode(min, max, amplitude, curve.ModeDefault|curve.ModeMinStart|curve.ModeMinEnd)
	if err != nil {
		return nil, err
	}
	return &Speed{
		min: min,
		max: max,
		gen: gen,
	}, nil
}

func (t *Speed) ToProto() *proto.TypeState {
	return &proto.TypeState{
		Min: t.min,
		Max: t.max,
		Val: t.val,
		Gen: t.gen.ToProto(),
	}
}

func (t *Speed) FromProto(speed *proto.TypeState) {
	t.gen = new(curve.Curve)
	t.gen.FromProto(speed.Gen)
	t.min = speed.Min
	t.max = speed.Max
	t.val = speed.Val
}

func (t *Speed) Min() float64 {
	return t.min
}

func (t *Speed) Max() float64 {
	return t.max
}

func (t *Speed) Value() float64 {
	return t.val
}

func (t *Speed) String() string {
	return fmt.Sprintf("speed: %.2f m/s", t.val)
}

func (t *Speed) Next(tick float64) {
	point := t.gen.Point(tick)
	t.val = point.Y
}
