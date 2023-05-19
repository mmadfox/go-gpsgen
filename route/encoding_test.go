package route

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncoding(t *testing.T) {
	routes, _ := RoutesForFrance()

	// encode
	data, err := Encode(routes)
	require.NoError(t, err)
	require.NotZero(t, data)

	// decode
	routes2, err := Decode(data)
	require.NoError(t, err)
	require.Equal(t, len(routes), len(routes2))

	for i := 0; i < len(routes); i++ {
		require.Equal(t, routes[i].Distance(), routes[i].Distance())
	}
}
