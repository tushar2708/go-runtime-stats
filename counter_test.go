package runtimestats

import (
	"testing"
)

func Test_perSecondCounter(t *testing.T) {
	type args struct {
		name  string
		value int64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "first counter should be zero",
			args: args{
				name:  "counterTest",
				value: 10,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := perSecondCounter(tt.args.name, tt.args.value); got != tt.want {
				t.Errorf("perSecondCounter() = %v, want %v", got, tt.want)
			}
		})
	}
}
