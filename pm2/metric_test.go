package pm2

import (
	"testing"

	a "github.com/stretchr/testify/assert"
)

func TestMetric(t *testing.T) {
	assert := a.New(t)
	ch := new(MockChannel)
	pm := getMockPm2(ch)
	assert.True(pm.isConnected())
	pm.prepareMetrics()
	pm.refreshMetricInfo()
	t.Logf("%+v", pm.metrics)
	assert.NotEmpty(pm.metrics)
	go pm.startSendStatus()
}
