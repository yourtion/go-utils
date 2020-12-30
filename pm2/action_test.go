package pm2

import (
	"testing"

	a "github.com/stretchr/testify/assert"
)

func TestAction(t *testing.T) {
	assert := a.New(t)
	ch := new(MockChannel)
	pm := getMockPm2(ch)
	assert.True(pm.isConnected())
	pm.AddAction("test", func(payload interface{}) string {
		return "ok"
	})
	_, ok := pm.actions.Load("test")
	assert.True(ok)
	ret := pm.callAction("test", nil)
	assert.Equal(*ret, "ok")

	pm.actionHandler("test", "ok")
}
