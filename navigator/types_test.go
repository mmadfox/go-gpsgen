package navigator

import (
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    Name
		wantErr bool
	}{
		{
			name: "should return error when name is too short",
			args: args{
				name: strings.Repeat("n", MinNameValue-1),
			},
			wantErr: true,
		},
		{
			name: "should return error when name is too long",
			args: args{
				name: strings.Repeat("n", MaxNameValue+1),
			},
			wantErr: true,
		},
		{
			name: "should return valid name when all params are valid",
			args: args{
				name: "somename",
			},
			want: Name{value: "somename"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseName(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseName() = %v, want %v", got, tt.want)
			}
			require.False(t, got.IsEmpty())
		})
	}
}
