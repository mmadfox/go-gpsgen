package curve

import (
	"errors"
	"math/rand"

	"github.com/mmadfox/go-gpsgen/proto"
)

var (
	ErrNoControlPoints = errors.New("curves: no control points")
)

type CurveMode uint16

const (
	ModeDefault = 1 << (16 - 1 - iota)
	ModeMinStart
	ModeMinEnd
)

type Point struct {
	X, Y float64
}

type point struct {
	vp Point
	cp Point
}

type Curve struct {
	points []point
	mode   int
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

func (c *Curve) ToProto() *proto.Curve {
	pbcurve := &proto.Curve{
		Points: make([]*proto.Curve_ControlPoint, len(c.points)),
		Mode:   int64(c.mode),
	}
	for i := 0; i < len(c.points); i++ {
		pbcurve.Points[i] = &proto.Curve_ControlPoint{
			Vp: &proto.Curve_Point{X: c.points[i].vp.X, Y: c.points[i].vp.Y},
			Cp: &proto.Curve_Point{X: c.points[i].cp.X, Y: c.points[i].cp.Y},
		}
	}
	return pbcurve
}

func (c *Curve) FromProto(pb *proto.Curve) {
	c.points = make([]point, len(pb.Points))
	for i := 0; i < len(pb.Points); i++ {
		pt := pb.Points[i]
		c.points[i] = point{
			vp: Point{
				X: pt.Vp.X,
				Y: pt.Vp.Y,
			},
			cp: Point{
				X: pt.Cp.X,
				Y: pt.Cp.Y,
			},
		}
	}
	c.mode = int(pb.Mode)
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
		mode:   int(m),
	}

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
