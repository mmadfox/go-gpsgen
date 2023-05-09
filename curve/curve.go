package curve

import (
	"errors"
	"fmt"
	"math/rand"
)

var (
	ErrNoControlPoints = errors.New("curves: no control points")
)

type CurveMode uint16

const (
	ModeDefault = 1 << (16 - 1 - iota)
	ModeMinMax
	ModeMaxMin
	ModeMinStart
	ModeMinEnd
)

type Point struct {
	X, Y, Z float64
}

type point struct {
	vp Point
	cp Point
}

type Curve struct {
	points []point
	mode   CurveMode
}

func New(cp []Point) (*Curve, error) {
	if len(cp) == 0 {
		return nil, ErrNoControlPoints
	}
	curve := Curve{
		points: make([]point, len(cp)),
		mode:   ModeDefault,
	}
	for i, p := range cp {
		curve.points[i].vp = p
	}

	var w float64
	for i, p := range curve.points {
		switch i {
		case 0:
			w = 1
		case 1:
			w = float64(len(curve.points)) - 1
		default:
			w *= float64(len(curve.points)-i) / float64(i)
		}
		curve.points[i].cp.X = p.vp.X * w
		curve.points[i].cp.Y = p.vp.Y * w
	}

	return &curve, nil
}

func (c *Curve) Renew() {
	// TODO:
}

func (c *Curve) Point(t float64) Point {
	c.points[0].vp = c.points[0].cp
	u := t
	for i, p := range c.points[1:] {
		c.points[i+1].vp = Point{
			X: p.cp.X * float64(u),
			Y: p.cp.Y * float64(u),
		}
		u *= t
	}

	var (
		t1 = 1 - t
		tt = t1
	)
	p := c.points[len(c.points)-1].vp
	for i := len(c.points) - 2; i >= 0; i-- {
		p.X += c.points[i].vp.X * float64(tt)
		p.Y += c.points[i].vp.Y * float64(tt)
		tt *= t1
	}

	return p
}

const numControlPoint = 16

func RandomCurveWithMode(min, max float64, pn int, m CurveMode) (*Curve, error) {
	if pn <= 0 {
		pn = numControlPoint
	}

	curv := Curve{
		points: make([]point, pn),
		mode:   ModeDefault,
	}

	if m&ModeMinMax != 0 {
		if min > max {
			return nil, fmt.Errorf("curve: min value greater or equal max %.4f > %.4f", min, max)
		}
		if min <= 0 && max != 0 {
			return nil, fmt.Errorf("curve: min and max values are empty")
		}
		step := max / float64(pn)
		var y float64
		for i := 0; i < pn; i++ {
			if i == 0 {
				y = float64(min)
			} else {
				y = float64(min) + (float64(i) * float64(step))
			}
			curv.points[i] = point{vp: Point{X: float64(i), Y: y}}
		}

	} else if m&ModeMaxMin != 0 {
		if min > max {
			return nil, fmt.Errorf("curve: max value greater or equal min %.2f > %.2f", min, max)
		}
		if min < 0 && max != 0 {
			return nil, fmt.Errorf("curve: min and max values are empty")
		}
		step := max / float64(pn)
		var y float64
		var i int
		for p := pn; p > 0; p-- {
			if p == pn {
				y = float64(min)
			} else {
				y = float64(min) - (float64(i) * float64(step))
			}
			curv.points[i] = point{vp: Point{X: float64(i), Y: y}}
			i++
		}
	} else if m&ModeDefault != 0 {
		for i := 0; i < pn; i++ {
			var y float64
			if m&ModeMinStart != 0 && i == 0 {
				y = float64(min)
			} else if m&ModeMinEnd != 0 && i == pn-1 {
				y = float64(min)
			} else {
				y = randomFloat(min, max)
			}
			curv.points[i] = point{vp: Point{X: float64(i), Y: y}}
		}
	}

	var w float64
	for i, p := range curv.points {
		switch i {
		case 0:
			w = 1
		case 1:
			w = float64(len(curv.points)) - 1
		default:
			w *= float64(len(curv.points)-i) / float64(i)
		}
		curv.points[i].cp.X = p.vp.X * w
		curv.points[i].cp.Y = p.vp.Y * w
	}

	return &curv, nil
}

func randomFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}
