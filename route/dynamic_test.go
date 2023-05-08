package route

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerate(t *testing.T) {
	for i := 0; i < 1000; i++ {
		route, err := Generate()
		require.NoError(t, err)
		require.NotNil(t, route)
		require.Equal(t, 1, route.NumTracks())
		require.GreaterOrEqual(t, 60, route.NumSegments(0))
	}
}

func TestGenerateFor(t *testing.T) {
	for _, country := range countries {
		for i := 0; i < 100; i++ {
			route, err := GenerateFor(country)
			require.NoError(t, err)
			require.NotNil(t, route)
			require.Equal(t, 1, route.NumTracks())
			require.GreaterOrEqual(t, 60, route.NumSegments(0))
		}
	}
}
