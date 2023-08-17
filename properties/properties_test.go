package properties

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProperties(t *testing.T) {
	tests := []struct {
		name    string
		arrange func(Properties)
		assert  func(Properties)
	}{
		{
			name: "should return string value",
			arrange: func(p Properties) {
				p.Set("key", "value")
			},
			assert: func(p Properties) {
				val, ok := p.String("key")
				require.True(t, ok)
				require.Equal(t, "value", val)
			},
		},
		{
			name: "should return int value",
			arrange: func(p Properties) {
				p.Set("key", 1000)
			},
			assert: func(p Properties) {
				val, ok := p.Int("key")
				require.True(t, ok)
				require.Equal(t, 1000, val)
			},
		},
		{
			name: "should return int when value is int32",
			arrange: func(p Properties) {
				p.Set("key", int32(1000))
			},
			assert: func(p Properties) {
				val, ok := p.Int("key")
				require.True(t, ok)
				require.Equal(t, 1000, val)
			},
		},
		{
			name: "should return int when value is int64",
			arrange: func(p Properties) {
				p.Set("key", int64(1000))
			},
			assert: func(p Properties) {
				val, ok := p.Int("key")
				require.True(t, ok)
				require.Equal(t, 1000, val)
			},
		},
		{
			name: "should return float64 value",
			arrange: func(p Properties) {
				p.Set("key", float64(1000))
			},
			assert: func(p Properties) {
				val, ok := p.Float64("key")
				require.True(t, ok)
				require.Equal(t, float64(1000), val)
			},
		},
		{
			name: "should return bool value",
			arrange: func(p Properties) {
				p.Set("key", true)
			},
			assert: func(p Properties) {
				val, ok := p.Bool("key")
				require.True(t, ok)
				require.True(t, val)
			},
		},
		{
			name: "should return merged value",
			arrange: func(p Properties) {
				p.Set("bool", true)
				p.Merge(Properties{"foo": "bar", "baz": 1})
			},
			assert: func(p Properties) {
				b, ok := p.Bool("bool")
				require.True(t, ok)
				require.True(t, b)
				f, ok := p.String("foo")
				require.True(t, ok)
				require.Equal(t, "bar", f)
				i, ok := p.Int("baz")
				require.True(t, ok)
				require.Equal(t, 1, i)
			},
		},
		{
			name: "should return default values",
			assert: func(p Properties) {
				v1, ok := p.Bool("")
				require.False(t, ok)
				require.False(t, v1)

				v2, ok := p.String("")
				require.False(t, ok)
				require.Empty(t, v2)

				v3, ok := p.Int("")
				require.False(t, ok)
				require.Zero(t, v3)

				v4, ok := p.Float64("")
				require.False(t, ok)
				require.Zero(t, v4)
			},
		},
		{
			name: "should removes items when properties is not empty",
			arrange: func(p Properties) {
				p.Set("key", "val")
			},
			assert: func(p Properties) {
				p.Remove("key")
				v, ok := p.String("key")
				require.False(t, ok)
				require.Empty(t, v)
			},
		},
		{
			name: "should resets when properties is not empty",
			arrange: func(p Properties) {
				p.Set("key", "val")
			},
			assert: func(p Properties) {
				p.Reset()
				require.Empty(t, p)
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			props := Make()
			if tc.arrange != nil {
				tc.arrange(props)
			}
			if tc.assert != nil {
				tc.assert(props)
			}
		})
	}
}
