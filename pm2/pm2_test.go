package pm2

import (
	"sync"
	"testing"

	a "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/yourtion/go-utils/nodejs"
)

type MockChannel struct {
	mock.Mock
}

func (m *MockChannel) Read() (*nodejs.NodeMessage, error) {
	return nil, nil
}

func (m *MockChannel) Write(*nodejs.NodeMessage) error {
	return nil
}

func getMockPm2(ch *MockChannel) *pm {
	pm := GetInstance()
	pm.connected = true
	pm.tran = &transporter{
		ch: ch,
		mu: sync.Mutex{},
	}
	return pm
}

func TestPm2(t *testing.T) {
	assert := a.New(t)
	ch := new(MockChannel)
	pm := getMockPm2(ch)
	assert.True(pm.isConnected())
}
