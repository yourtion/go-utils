package pm2

import (
	"encoding/json"
	"sync"

	"github.com/yourtion/go-utils/nodejs"
)

type transporter struct {
	ch nodejs.NodeChannel
	mu sync.Mutex
}

type messageHandler func(msg string, opts string)

// connect 连接上 pm2 父进程
func connect() (*transporter, error) {
	channel, err := nodejs.RunAsNodeChild()
	if err != nil {
		return nil, err
	}
	t := &transporter{
		ch: channel,
		mu: sync.Mutex{},
	}
	return t, nil
}

// setHandler 配置并开启消息接收，注意使用 goroutine 启动
func (t *transporter) setHandler(handler messageHandler) {
	for {
		msg, err := t.ch.Read()
		if err != nil {
			continue
		}
		act := new(actionMessage)
		err = json.Unmarshal(msg.Message, &act)
		if err != nil {
			err = json.Unmarshal(msg.Message, &act.Msg)
			if err != nil {
				continue
			}
		}
		handler(act.Msg, act.Opts)
	}
}

// sendJson 给 pm2 发送消息
func (t *transporter) sendJson(msg interface{}) {
	b, err := json.Marshal(msg)
	// fmt.Printf("%s", b)
	if err != nil {
		return
	}

	t.mu.Lock()
	defer t.mu.Unlock()
	err = t.ch.Write(&nodejs.NodeMessage{Message: b})
}
