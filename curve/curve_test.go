package curve

import (
	"testing"

	"github.com/mmadfox/go-gpsgen/proto"
	"github.com/stretchr/testify/require"
)

func TestCurveToProto(t *testing.T) {
	controlPoint := 16
	curve, err := RandomCurveWithMode(1, 10, controlPoint, ModeDefault)
	require.NoError(t, err)
	require.NotNil(t, curve)

	protoCurve := curve.ToProto()
	require.NotNil(t, protoCurve)
	require.Equal(t, controlPoint, len(protoCurve.Points))
	require.Equal(t, int(ModeDefault), int(protoCurve.Mode))
	for i := 0; i < len(curve.points); i++ {
		require.Equal(t, curve.points[i].cp.X, protoCurve.Points[i].Cp.X)
		require.Equal(t, curve.points[i].cp.Y, protoCurve.Points[i].Cp.Y)
		require.Equal(t, curve.points[i].vp.X, protoCurve.Points[i].Vp.X)
		require.Equal(t, curve.points[i].vp.Y, protoCurve.Points[i].Vp.Y)
	}
}

func TestCurveFromProto(t *testing.T) {
	protoCurve := &proto.Curve{
		Points: []*proto.Curve_ControlPoint{
			{
				Vp: &proto.Curve_Point{X: 1, Y: 1},
				Cp: &proto.Curve_Point{X: 1, Y: 1},
			},
			{
				Vp: &proto.Curve_Point{X: 5, Y: 10},
				Cp: &proto.Curve_Point{X: 5, Y: 10},
			},
		},
		Mode: int64(ModeMinEnd),
	}

	curv := new(Curve)
	curv.FromProto(protoCurve)
	require.Equal(t, len(protoCurve.Points), len(curv.points))
	for i := 0; i < len(protoCurve.Points); i++ {
		require.Equal(t, protoCurve.Points[i].Cp.X, curv.points[i].cp.X)
		require.Equal(t, protoCurve.Points[i].Cp.Y, curv.points[i].cp.Y)
		require.Equal(t, protoCurve.Points[i].Vp.X, curv.points[i].vp.X)
		require.Equal(t, protoCurve.Points[i].Vp.Y, curv.points[i].vp.Y)
	}
	require.Equal(t, int(protoCurve.Mode), curv.mode)
}
