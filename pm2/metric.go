package pm2

import (
	"time"

	"github.com/yourtion/go-utils/metric"
)

// AddMetric to global metrics array
func (pm2 *pm) AddMetric(metric *metric.Metric) {
	if !pm2.connected {
		return
	}
	pm2.metricLock.Lock()
	defer pm2.metricLock.Unlock()
	pm2.metrics[metric.Name] = metric
}

func (pm2 *pm) refreshMetricInfo() {
	if !pm2.connected {
		return
	}
	pm2.metricLock.RLock()
	defer pm2.metricLock.RUnlock()
	if pm2.metricHandler != nil {
		(*pm2.metricHandler)()
	}
	for _, m := range pm2.metrics {
		m.Get()
	}
}

func (pm2 *pm) attachMetricHandler(handler func()) {
	pm2.metricHandler = &handler
}

func (pm2 *pm) prepareMetrics() {
	initMetricsMemStats()
	pm2.attachMetricHandler(memStatsHandler)
	pm2.AddMetric(goRoutines())
	pm2.AddMetric(cgoCalls())
	pm2.AddMetric(globalMetricsMemStats.NumGC)
	pm2.AddMetric(globalMetricsMemStats.NumMallocs)
	pm2.AddMetric(globalMetricsMemStats.NumFree)
	pm2.AddMetric(globalMetricsMemStats.HeapAlloc)
	pm2.AddMetric(globalMetricsMemStats.Pause)
}

func (pm2 *pm) startSendStatus() {
	ticker := time.NewTicker(time.Second)
	for {
		<-ticker.C
		pm2.refreshMetricInfo()
		pm2.tran.sendJson(message{Type: "axm:monitor", Data: pm2.metrics})
	}
}
