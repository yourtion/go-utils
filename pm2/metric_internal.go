package pm2

import (
	"runtime"

	"github.com/yourtion/go-utils/metric"
)

// MetricsMemStats is a structure to simplify storage of mem values
type MetricsMemStats struct {
	Init      bool
	NumGC     *metric.Metric
	LastNumGC float64

	NumMallocs     *metric.Metric
	LastNumMallocs float64

	NumFree     *metric.Metric
	LastNumFree float64

	HeapAlloc *metric.Metric

	Pause     *metric.Metric
	LastPause float64
}

// GlobalMetricsMemStats store current and last mem stats
var globalMetricsMemStats MetricsMemStats

// GoRoutines create a func metric who return number of current GoRoutines
func goRoutines() *metric.Metric {
	return metric.CreateFuncMetric("GoRoutines", "metric", "routines", func() float64 {
		return float64(runtime.NumGoroutine())
	})
}

// CgoCalls create a func metric who return number of current C calls of last second
func cgoCalls() *metric.Metric {
	last := runtime.NumCgoCall()
	return metric.CreateFuncMetric("CgoCalls/sec", "metric", "calls/sec", func() float64 {
		calls := runtime.NumCgoCall()
		v := calls - last
		last = calls
		return float64(v)
	})
}

// InitMetricsMemStats create metrics for MemStats
func initMetricsMemStats() {
	numGC := metric.CreateMetric("GCRuns/sec", "metric", "runs")
	numMalloc := metric.CreateMetric("mallocs/sec", "metric", "mallocs")
	numFree := metric.CreateMetric("free/sec", "metric", "frees")
	heapAlloc := metric.CreateMetric("heapAlloc", "metric", "bytes")
	pause := metric.CreateMetric("Pause/sec", "metric", "ns/sec")

	globalMetricsMemStats = MetricsMemStats{
		Init:       true,
		NumGC:      numGC,
		NumMallocs: numMalloc,
		NumFree:    numFree,
		HeapAlloc:  heapAlloc,
		Pause:      pause,
	}
}

// Handler write values in MemStats metrics
func memStatsHandler() {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)

	globalMetricsMemStats.NumGC.Set(float64(stats.NumGC) - globalMetricsMemStats.LastNumGC)
	globalMetricsMemStats.LastNumGC = float64(stats.NumGC)

	globalMetricsMemStats.NumMallocs.Set(float64(stats.Mallocs) - globalMetricsMemStats.LastNumMallocs)
	globalMetricsMemStats.LastNumMallocs = float64(stats.Mallocs)

	globalMetricsMemStats.NumFree.Set(float64(stats.Frees) - globalMetricsMemStats.LastNumFree)
	globalMetricsMemStats.LastNumFree = float64(stats.Frees)

	globalMetricsMemStats.HeapAlloc.Set(float64(stats.HeapAlloc))

	globalMetricsMemStats.Pause.Set(float64(stats.PauseTotalNs) - globalMetricsMemStats.LastPause)
	globalMetricsMemStats.LastPause = float64(stats.PauseTotalNs)
}
