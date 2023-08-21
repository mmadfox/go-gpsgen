package random

import (
	"testing"
)

func TestPolygon(t *testing.T) {
	type args struct {
		points int
		zoom   float64
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "should return  defaultPoints when num of points less or then 0",
			args: args{
				points: -1,
			},
			want: defaultPoints + 1,
		},
		{
			name: "should generate valid points when zoom < 0",
			args: args{
				points: 32,
				zoom:   -1.0,
			},
			want: 32 + 1,
		},
		{
			name: "should generate valid points when zoom > 1000",
			args: args{
				points: 32,
				zoom:   2000,
			},
			want: 32 + 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Polygon(tt.args.points, tt.args.zoom)
			if len(got) != tt.want {
				t.Errorf("Polygon() = %v, want %v number of points", got, tt.want)
			}
		})
	}
}
