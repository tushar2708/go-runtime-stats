package runtimestats

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_sanitizeMetricName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1",
			args: args{
				name: "mymetricname",
			},
			want: "mymetricname",
		},
		{
			name: "2",
			args: args{
				name: "my metric name",
			},
			want: "my_metric_name",
		},
		{
			name: "3",
			args: args{
				name: "my/metric/name",
			},
			want: "my_metric_name",
		},
		{
			name: "4",
			args: args{
				name: "my.metric name",
			},
			want: "my_metric_name",
		},
		{
			name: "5",
			args: args{
				name: "my-metric/name",
			},
			want: "my-metric_name",
		},
		{
			name: "6",
			args: args{
				name: "my-metric@name",
			},
			want: "my-metricname",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sanitizeMetricName(tt.args.name); got != tt.want {
				t.Errorf("sanitizeMetricName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newRuntimeStats(t *testing.T) {
	type args struct {
		statsDPrefix string
		tags         []string
	}
	testsPrefix := []struct {
		name string
		args args
		want string
	}{
		{
			name: "default statsd prefix should be set",
			args: args{
				statsDPrefix: "app.app_env",
				tags:         []string{},
			},
			want: "app.app_env",
		},
	}
	for _, tt := range testsPrefix {
		t.Run(tt.name, func(t *testing.T) {
			if got := newRuntimeStats(tt.args.statsDPrefix, tt.args.tags...); !(got.StatsDPrefix == tt.want) {
				t.Errorf("newRuntimeStats() = %v, want %v", got, tt.want)
			}
		})
	}

	testsBase := []struct {
		name string
		args args
		want string
	}{
		{
			name: "matric base should be correct",
			args: args{
				statsDPrefix: "app.app_env",
				tags:         []string{},
			},
			want: "app.app_env.runtime.",
		},
	}
	for _, tt := range testsBase {
		t.Run(tt.name, func(t *testing.T) {
			if got := newRuntimeStats(tt.args.statsDPrefix, tt.args.tags...); !(got.getMetricBase() == tt.want) {
				t.Errorf("newRuntimeStats() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Start(t *testing.T) {

	type args struct {
		statsdHost   string
		statsDPrefix string
		tags         []string
		interval     int
	}
	type want struct {
		err        error
		matricBase string
		hostName   string
		interval   time.Duration
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "default statsd prefix should be set",
			args: args{
				statsdHost:   "0.0.0.0:8015",
				statsDPrefix: "app.app_env",
				tags:         []string{},
				interval:     5,
			},
			want: want{
				err:        nil,
				matricBase: "app.app_env.runtime.",
				hostName:   "0.0.0.0:8015",
				interval:   time.Duration(5 * time.Second),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, gotErr := Start(tt.args.statsdHost, tt.args.statsDPrefix, tt.args.interval)

			if !(gotErr == nil) {
				t.Errorf("gotErr = %v, want = nil", gotErr)
			}

			if !(got.getMetricBase() == tt.want.matricBase) {
				t.Errorf("got.getMetricBase() = %v, tt.want.matricBase = %v", got.getMetricBase(), tt.want.matricBase)
			}

			if !(got.StatsdHost == tt.want.hostName) {
				t.Errorf("got.StatsdHost = %v, tt.want.hostName = %v", got.StatsdHost, tt.want.hostName)
			}

			if !(got.PublishInterval == tt.want.interval) {
				t.Errorf("got.PublishInterval = %v, tt.want.interval = %v", got.PublishInterval, tt.want.interval)
			}
		})
	}
}

func returnsKeys(t *testing.T, expectedKeys []string, response map[string]float64) {
	for _, k := range expectedKeys {
		_, found := response[k]
		assert.True(t, found, fmt.Sprintf("Should expose metric %s", k))
	}
}

func Test_MemStats(t *testing.T) {
	returnsKeys(t, []string{
		"memory.objects.HeapObjects",
		"memory.summary.Alloc",
		"memory.counters.Mallocs",
		"memory.counters.Frees",
		"memory.summary.System",
		"memory.heap.Idle",
		"memory.heap.InUse",
		"memory.stack.InUse",
	}, memStats())
}

func Test_GoRoutines(t *testing.T) {
	returnsKeys(t, []string{
		"goroutines.total",
	}, goRoutines())
}

func Test_CgoCalls(t *testing.T) {
	returnsKeys(t, []string{
		"cgo.calls",
	}, cgoCalls())
}

func Test_Gcs(t *testing.T) {
	gcs := gcs()
	returnsKeys(t, []string{
		"gc.perSecond",
		"gc.pauseTimeNs",
		"gc.pauseTimeMs",
	}, gcs)

	assert.Equal(t, gcs["gc.pauseTimeNs"]/float64(1000000), gcs["gc.pauseTimeMs"], "Pause time NS should convert to MS")
}
