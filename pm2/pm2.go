package pm2

import (
	"os"
	"sync"
)

const pm2InsKey = "PM2_INTERACTOR_PROCESSING"

// pm2 实例
type pm struct {
	connected bool         // 是否连接
	tran      *transporter // 传输
	actions   sync.Map     // 全局 actions
}

var instance *pm
var once sync.Once

// GetInstance 获取 pm2 单例
func GetInstance() *pm {
	once.Do(func() {
		instance = create()
	})
	return instance
}

// create 创建 pm2 实例
func create() *pm {
	pm2 := &pm{
		connected: false,
		tran:      nil,
		actions:   sync.Map{},
	}
	if os.Getenv(pm2InsKey) != "true" {
		return pm2
	}
	t, err := connect()
	if err != nil {
		return pm2
	}
	pm2.connected = true
	pm2.tran = t
	go pm2.tran.setHandler(pm2.actionHandler)
	return pm2
}
