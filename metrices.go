package runtimestats

import "runtime"

func memStats() map[string]float64 {
	m := runtime.MemStats{}
	runtime.ReadMemStats(&m)
	metrics := map[string]float64{
		"memory.objects.HeapObjects": float64(m.HeapObjects),
		"memory.summary.Alloc":       float64(m.Alloc),
		"memory.counters.Mallocs":    perSecondCounter("mallocs", int64(m.Mallocs)),
		"memory.counters.Frees":      perSecondCounter("frees", int64(m.Frees)),
		"memory.summary.System":      float64(m.HeapSys),
		"memory.heap.Idle":           float64(m.HeapIdle),
		"memory.heap.InUse":          float64(m.HeapInuse),
		"memory.stack.InUse":         float64(m.StackInuse),
	}

	return metrics
}

func goRoutines() map[string]float64 {
	return map[string]float64{
		"goroutines.total": float64(runtime.NumGoroutine()),
	}
}

var lastGcPause float64
var lastGcTime uint64
var lastGcPeriod float64

func gcs() map[string]float64 {
	m := runtime.MemStats{}
	runtime.ReadMemStats(&m)
	gcPause := float64(m.PauseNs[(m.NumGC+255)%256])
	if gcPause > 0 {
		lastGcPause = gcPause
	}

	if m.LastGC > lastGcTime {
		lastGcPeriod = float64(m.LastGC - lastGcTime)
		if lastGcPeriod == float64(m.LastGC) {
			lastGcPeriod = 0
		}

		lastGcPeriod = lastGcPeriod / 1000000

		lastGcTime = m.LastGC
	}

	return map[string]float64{
		"gc.perSecond":   perSecondCounter("gcs-total", int64(m.NumGC)),
		"gc.pauseTimeNs": lastGcPause,
		"gc.pauseTimeMs": lastGcPause / float64(1000000),
		"gc.period":      lastGcPeriod,
	}
}

func cgoCalls() map[string]float64 {
	return map[string]float64{
		"cgo.calls": perSecondCounter("cgoCalls", runtime.NumCgoCall()),
	}
}
