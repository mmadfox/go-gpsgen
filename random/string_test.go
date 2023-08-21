package random

import (
	"testing"
)

func TestString(t *testing.T) {
	type args struct {
		length int
	}
	tests := []struct {
		name    string
		args    args
		wantLen int
	}{
		{
			name: "should return 0 chars when lenght 0",
			args: args{
				length: 0,
			},
			wantLen: 0,
		},
		{
			name: "should return 0 char when lenght -1",
			args: args{
				length: -1,
			},
			wantLen: 0,
		},
		{
			name: "should return 4 char when lenght 4",
			args: args{
				length: 4,
			},
			wantLen: 4,
		},
		{
			name: "should return 32 char when lenght 32",
			args: args{
				length: 32,
			},
			wantLen: 32,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := String(tt.args.length)
			if len(got) != tt.wantLen {
				t.Fatalf("String() => %s(%d), want %d lenght", got, len(got), tt.wantLen)
			}
		})
	}
}
