package curve

import (
	"sort"
	"testing"

	"github.com/mmadfox/go-gpsgen/proto"
	"github.com/stretchr/testify/require"
)

func TestCurve_New(t *testing.T) {
	type args struct {
		min  float64
		max  float64
		cp   int
		mode CurveMode
	}
	tests := []struct {
		name   string
		args   args
		assert func(c *Curve)
	}{
		{
			name: "should return curve with default control points",
			args: args{
				cp: 0,
			},
			assert: func(c *Curve) {
				require.Equal(t, defaultControlPoint, c.NumControlPoints())
			},
		},
		{
			name: "should return curve with default control points when negative cp",
			args: args{
				cp: -1,
			},
			assert: func(c *Curve) {
				require.Equal(t, defaultControlPoint, c.NumControlPoints())
			},
		},
		{
			name: "should return curve when mode minEnd",
			args: args{
				min:  1,
				max:  10,
				mode: ModeMinEnd,
			},
			assert: func(c *Curve) {
				lastCP := c.points[len(c.points)-1]
				require.Equal(t, 1.0, lastCP.vp.Y)
			},
		},
		{
			name: "should return curve when mode ",
			args: args{
				min:  2,
				max:  10,
				mode: ModeMinStart,
			},
			assert: func(c *Curve) {
				firstCP := c.points[0]
				require.Equal(t, 2.0, firstCP.vp.Y)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			curve := New(tt.args.min, tt.args.max, tt.args.cp, tt.args.mode)
			if tt.assert != nil {
				tt.assert(curve)
			}
		})
	}
}

func TestCurve_Point(t *testing.T) {
	type args struct {
		min  float64
		max  float64
		cp   int
		mode CurveMode
	}
	tests := []struct {
		name   string
		args   args
		assert func(c *Curve)
	}{
		{
			name: "should return point with positive values",
			args: args{
				min: 1,
				max: 2,
				cp:  8,
			},
			assert: func(c *Curve) {
				values := make([]float64, 100)
				for i := 1; i <= 100; i++ {
					T := float64(i) * 0.01
					pt := c.Point(T)
					require.NotZero(t, pt.Y)
					require.NotZero(t, pt.X)
					values[i-1] = pt.Y
				}
				sort.Float64s(values)
				require.GreaterOrEqual(t, values[0], 1.0)
				require.LessOrEqual(t, values[len(values)-1], 2.0)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			curve := New(tt.args.min, tt.args.max, tt.args.cp, tt.args.mode)
			if tt.assert != nil {
				tt.assert(curve)
			}
		})
	}
}

func TestCurve_Snapshot(t *testing.T) {
	type args struct {
		min  float64
		max  float64
		cp   int
		mode CurveMode
	}
	tests := []struct {
		name   string
		args   args
		assert func(c *Curve, s *proto.Snapshot_Curve)
	}{
		{
			name: "should return snapshot when all params are valid",
			args: args{
				min:  2,
				max:  400,
				cp:   16,
				mode: ModeMinStart,
			},
			assert: func(c *Curve, s *proto.Snapshot_Curve) {
				require.Equal(t, c.Min(), s.Min)
				require.Equal(t, c.Max(), s.Max)
				require.Equal(t, c.mode, int(s.Mode))
				require.Len(t, s.Points, c.NumControlPoints())

				for i := 0; i < c.NumControlPoints(); i++ {
					p1 := c.points[i]
					p2 := s.Points[i]
					require.Equal(t, p1.vp.Y, p2.Vp.Y)
					require.Equal(t, p1.vp.X, p2.Vp.X)
					require.Equal(t, p1.cp.Y, p2.Cp.Y)
					require.Equal(t, p1.cp.X, p2.Cp.X)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			curve := New(tt.args.min, tt.args.max, tt.args.cp, tt.args.mode)
			if tt.assert != nil {
				tt.assert(curve, curve.Snapshot())
			}
		})
	}
}

func TestCurve_FromSnapshot(t *testing.T) {
	curve := New(1, 100, 8, ModeDefault)
	snap := curve.Snapshot()
	curve2 := new(Curve)
	curve2.FromSnapshot(snap)
	require.Equal(t, curve.NumControlPoints(), curve2.NumControlPoints())
	require.Equal(t, curve.Min(), curve2.Min())
	require.Equal(t, curve.Max(), curve2.Max())
	require.Equal(t, curve.mode, curve2.mode)

	for i := 0; i < curve.NumControlPoints(); i++ {
		p1 := curve.points[i]
		p2 := curve2.points[i]
		require.Equal(t, p1.vp.Y, p2.vp.Y)
		require.Equal(t, p1.vp.X, p2.vp.X)
		require.Equal(t, p1.cp.Y, p2.cp.Y)
		require.Equal(t, p1.cp.X, p2.cp.X)
	}
}

func TestCurve_Shuffle(t *testing.T) {
	expectedControlPoints := 16
	c1 := New(1, 100, expectedControlPoints, ModeDefault)
	require.NotNil(t, c1)
	require.Len(t, c1.points, expectedControlPoints)

	diff := 0
	for l := 0; l < 3; l++ {
		prevControlPoints := make([]point, len(c1.points))
		copy(prevControlPoints, c1.points)
		c1.Shuffle()
		require.Len(t, c1.points, expectedControlPoints)
		for i := 0; i < expectedControlPoints; i++ {
			p1 := prevControlPoints[i]
			p2 := c1.points[i]
			if p1.vp.Y != p2.vp.Y {
				diff++
			}
		}
	}

	total := 3 * expectedControlPoints
	if diff != total {
		t.Fatalf("Shuffle() => %d, want > %d", diff, total)
	}
}
