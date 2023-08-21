package curve

import (
	"github.com/mmadfox/go-gpsgen/proto"
	"github.com/mmadfox/go-gpsgen/random"
)

// CurveMode represents different curve modes.
type CurveMode uint16

const (
	ModeDefault = 1 << (16 - 1 - iota)
	ModeMinStart
	ModeMinEnd
)

var rnd = random.NewRandom()

// Point represents a point in 2D space.
type Point struct {
	X, Y float64
}

type point struct {
	vp Point
	cp Point
}

// Curve represents a curve with control points and related operations.
type Curve struct {
	points   []point
	mode     int
	min, max float64
}

// Min returns the minimum value of the curve.
func (c *Curve) Min() float64 {
	return c.min
}

// Max returns the maximum value of the curve.
func (c *Curve) Max() float64 {
	return c.max
}

// NumControlPoints returns the number of control points in the curve.
func (c *Curve) NumControlPoints() int {
	return len(c.points)
}

// Shuffle randomly shuffles the control points of the curve.s
func (c *Curve) Shuffle() {
	for i := 0; i < len(c.points); i++ {
		var y float64
		if c.mode&ModeMinStart != 0 && i == 0 {
			y = float64(c.min)
		} else if c.mode&ModeMinEnd != 0 && i == len(c.points)-1 {
			y = float64(c.min)
		} else {
			y = randomFloat(c.min, c.max)
		}
		c.points[i] = point{vp: Point{X: float64(i), Y: y}}
	}

	var w float64
	for i, p := range c.points {
		switch i {
		case 0:
			w = 1
		case 1:
			w = float64(len(c.points)) - 1
		default:
			w *= float64(len(c.points)-i) / float64(i)
		}
		c.points[i].cp.X = p.vp.X * w
		c.points[i].cp.Y = p.vp.Y * w
	}
}

// Point returns a point on the curve for a given parameter t.
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

// Snapshot generates a snapshot of the curve's state.
func (c *Curve) Snapshot() *proto.Snapshot_Curve {
	points := make([]*proto.Snapshot_Curve_ControlPoint, len(c.points))
	for i := 0; i < len(c.points); i++ {
		point := c.points[i]
		points[i] = &proto.Snapshot_Curve_ControlPoint{
			Vp: &proto.Snapshot_Curve_Point{
				X: point.vp.X,
				Y: point.vp.Y,
			},
			Cp: &proto.Snapshot_Curve_Point{
				X: point.cp.X,
				Y: point.cp.Y,
			},
		}
	}
	return &proto.Snapshot_Curve{
		Mode:   int64(c.mode),
		Min:    c.min,
		Max:    c.max,
		Points: points,
	}
}

// FromSnapshot restores the curve's state from a snapshot.
func (c *Curve) FromSnapshot(snap *proto.Snapshot_Curve) {
	c.mode = int(snap.Mode)
	c.min = snap.Min
	c.max = snap.Max
	c.points = make([]point, len(snap.Points))
	for i := 0; i < len(snap.Points); i++ {
		c.points[i] = point{
			vp: Point{X: snap.Points[i].Vp.X, Y: snap.Points[i].Vp.Y},
			cp: Point{X: snap.Points[i].Cp.X, Y: snap.Points[i].Cp.Y},
		}
	}
}

const defaultControlPoint = 4

// New creates a new Curve with specified parameters.
func New(min, max float64, controlPoints int, m CurveMode) *Curve {
	if controlPoints <= 0 {
		controlPoints = defaultControlPoint
	}

	curv := Curve{
		points: make([]point, controlPoints),
		mode:   int(m),
	}

	for i := 0; i < controlPoints; i++ {
		var y float64
		if m&ModeMinStart != 0 && i == 0 {
			y = float64(min)
		} else if m&ModeMinEnd != 0 && i == controlPoints-1 {
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

	curv.min = min
	curv.max = max

	return &curv
}

func randomFloat(min, max float64) float64 {
	return min + rnd.Float64()*(max-min)
}
