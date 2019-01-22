package runtimestats

import (
	"fmt"
	"gopkg.in/alexcesaro/statsd.v2"
	"regexp"
	"strings"
	"time"
)

type metrices []func() map[string]float64

//RuntimeStats stores thefields requied for stats publihing
type RuntimeStats struct {
	StatsdHost   string
	StatsDPrefix string

	MetricesToColllect metrices
	statsdClient       *statsd.Client

	PublishInterval time.Duration
	PublishTicker   *time.Ticker

	tags []string
}

func sanitizeMetricName(name string) string {
	for _, c := range []string{"/", ".", " "} {
		name = strings.Replace(name, c, "_", -1)
	}

	r := regexp.MustCompile("[^a-zA-Z0-9-_]")
	name = r.ReplaceAllString(name, "")

	return name
}

func newRuntimeStats(statsDPrefix string, tags ...string) *RuntimeStats {
	s := RuntimeStats{}
	s.StatsDPrefix = statsDPrefix

	s.MetricesToColllect = metrices{memStats, goRoutines, cgoCalls, gcs}
	s.tags = tags
	return &s
}

// Start creates a new insnce of NewRuntimeStats andarts themonitoring go-routine to pulish changing stats
func Start(statsdHost string, statsDPrefix string, publishInterval int, tags ...string) (*RuntimeStats, error) {
	s := newRuntimeStats(statsDPrefix, tags...)

	s.StatsdHost = statsdHost
	s.PublishInterval = time.Duration(publishInterval) * time.Second

	err := s.start()

	return s, err
}

func (s *RuntimeStats) getMetricBase() string {
	return strings.Join([]string{s.StatsDPrefix, "runtime", ""}, ".")
}

func (s *RuntimeStats) start() error {
	var err error
	s.statsdClient, err = statsd.New(statsd.Prefix(s.getMetricBase()), statsd.TagsFormat(statsd.InfluxDB),
		statsd.Tags(s.tags...), statsd.Address(s.StatsdHost))
	if err != nil {
		return err
	}
	s.PublishTicker = time.NewTicker(s.PublishInterval)

	// start the actual goroutine to send stats
	go s.startStatsPolling()

	return nil
}

func (s *RuntimeStats) startStatsPolling() {
	for {
		select {
		case <-s.PublishTicker.C:
			s.doSend()
		}
	}
	fmt.Println("done")
}

func (s *RuntimeStats) doSend() {
	for _, metrixFunc := range s.MetricesToColllect {
		metrics := metrixFunc()

		for metricName, metricValue := range metrics {
			s.statsdClient.Gauge(metricName, metricValue)
		}
	}
}

// Stop stops the run-time stats monitor
func (s *RuntimeStats) Stop() {
	s.PublishTicker.Stop()
}
