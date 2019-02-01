# go-runtime-stats

A tiny library to publish Golang runtime stats to Grafana (with or without InfluxDB tags)


[![GitHub license](https://img.shields.io/github/license/mashape/apistatus.svg)]()

[![Build Status](https://travis-ci.com/tushar2708/go-runtime-stats.svg?branch=master)](https://travis-ci.com/tushar2708/go-runtime-stats)

[![Maintainability](https://api.codeclimate.com/v1/badges/9eeb062d61505334f23b/maintainability)](https://codeclimate.com/github/tushar2708/go-runtime-stats/maintainability)

[![Test Coverage](https://api.codeclimate.com/v1/badges/9eeb062d61505334f23b/test_coverage)](https://codeclimate.com/github/tushar2708/go-runtime-stats/test_coverage)

[![Known Vulnerabilities](https://snyk.io/test/github/tushar2708/go-runtime-stats/badge.svg)](https://snyk.io/test/github/tushar2708/go-runtime-stats) 


Usage;

```go
package main

import "github.com/tushar2708/go-runtime-stats"

func main(){
	runtimestats.Start("statsd-host:8125", "app_name.app_env", 5, "tag1", "value1", "tag2", "value2")
}
```

Metrics exported;

| Metric                     | Source                           | Description                            | Unit               |
|----------------------------|----------------------------------|----------------------------------------|--------------------|
| cgo.calls                  | runtime.NumCgoCall()             | Number of Cgo Calls                    | calls per second   |
| gc.pauseTimeMs             | runtime.ReadMemStats             | Pause time of last GC run              | MS                 |
| gc.pauseTimeNs             | runtime.ReadMemStats             | Pause time of last GC run              | NS                 |
| gc.period                  | runtime.ReadMemStats             | Time between last two GC runs          | MS                 |
| gc.perSecond               | runtime.ReadMemStats             | Number of GCs per second               | runs per second    |
| goroutines.total           | runtime.NumGoroutine()           | Number of currently running goroutines | total              |
| memory.counters.Frees      | runtime.ReadMemStats.Frees       | Number of frees issued to the system   | frees per second   |
| memory.counters.Mallocs    | runtime.ReadMemStats.Mallocs     | Number of Mallocs issued to the system | mallocs per second |
| memory.heap.Idle           | runtime.ReadMemStats.HeapIdle    | Memory on the heap not in use          | bytes              |
| memory.heap.InUse          | runtime.ReadMemStats.HeapInuse   | Memory on the heap in use              | bytes              |
| memory.objects.HeapObjects | runtime.ReadMemStats.HeapObjects | Total objects on the heap              | # Objects          |
| memory.summary.Alloc       | runtime.ReadMemStats.Alloc       | Total bytes allocated                  | bytes              |
| memory.summary.System      | runtime.ReadMemStats.HeapSys     | Total bytes acquired from system       | bytes              |
