package pm2

// message 消息基础结构
type message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// actionFun 注册 action 具体的 func
type actionFun func(payload interface{}) string

// action 操作结构体
type action struct {
	ActionName string    `json:"action_name"` // action 名称
	ActionType string    `json:"action_type"` // default: "custom" else "internal"
	Callback   actionFun `json:"-"`           // 执行方法
}

// actionMessage 操作触发消息结构
type actionMessage struct {
	Name string `json:"name"` // 进程名称
	Msg  string `json:"msg"`  // 方法名称
	Opts string `json:"opts"` // 业务参数
}
