package route

import (
	"testing"

	"github.com/mmadfox/go-gpsgen/navigator"
	"github.com/stretchr/testify/require"
)

func TestStatic(t *testing.T) {
	routes := []func() (*navigator.Route, error){
		Russia1,
		Russia2,
		Russia3,
		Russia4,
		Russia5,
		France1,
		France2,
		France3,
		France4,
		France5,
		Spain1,
		Spain2,
		Spain3,
		Spain4,
		Spain5,
		China1,
		China2,
		China3,
		China4,
		China5,
	}
	for _, fn := range routes {
		r, err := fn()
		require.NoError(t, err)
		require.NotNil(t, r)
	}
}
