package pm2

import (
	"os"
	"sync"

	"github.com/yourtion/go-utils/metric"
)

const pm2HomeKey = "PM2_HOME"

type logFunc func(format string, args ...interface{})

func logNon(format string, args ...interface{}) {}

var log logFunc

// pm2 实例
type pm struct {
	connected     bool                      // 是否连接
	tran          *transporter              // 传输
	actions       sync.Map                  // 全局 actions
	metrics       map[string]*metric.Metric // 全局 metrics
	metricLock    sync.RWMutex              // metrics 读写锁
	metricHandler *func()                   // metrics 更新函数
}

var instance *pm
var once sync.Once

// GetInstance 获取 pm2 单例
func GetInstance() *pm {
	return GetInstanceWithLogger(logNon)
}

func GetInstanceWithLogger(logger logFunc) *pm {
	once.Do(func() {
		log = logger
		instance = create()
	})
	return instance
}

func (pm2 *pm) isConnected() bool {
	return pm2 != nil && pm2.connected && pm2.tran != nil
}

// create 创建 pm2 实例
func create() *pm {
	var pm2 = &pm{
		connected:     false,
		tran:          nil,
		actions:       sync.Map{},
		metrics:       make(map[string]*metric.Metric),
		metricLock:    sync.RWMutex{},
		metricHandler: nil,
	}
	log("ENV: %s\n", os.Getenv(pm2HomeKey))
	if os.Getenv(pm2HomeKey) == "" {
		return pm2
	}
	t, err := connect()
	if err != nil {
		log("connect error: %s\n", err)
		return pm2
	}
	pm2.connected = true
	pm2.tran = t
	pm2.prepareMetrics()
	go pm2.tran.setHandler(pm2.actionHandler)
	go pm2.startSendStatus()
	log("prepare done\n")
	return pm2
}
