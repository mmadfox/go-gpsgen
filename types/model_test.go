package types

import (
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewModel(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		args    args
		want    Model
		wantErr bool
	}{
		{
			name: "should return error, when model value too short",
			args: args{
				value: "",
			},
			wantErr: true,
		},
		{
			name: "should return error, when model value too long",
			args: args{
				value: strings.Repeat("a", maxModelLen+1),
			},
			wantErr: true,
		},
		{
			name: "should return valid model",
			args: args{value: "model"},
			want: Model{val: "model"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewModel(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewModel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewModel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModelRandom(t *testing.T) {
	model := RandomModel()
	require.False(t, model.IsEmpty())
	require.Len(t, model.String(), 9)
}
