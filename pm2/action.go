package pm2

// addAction 添加 action 到全局 map 中
func (pm2 *pm) addAction(action *action) {
	pm2.actions.Store(action.ActionName, action)
}

// callAction 通过 name 执行对应的 action
func (pm2 *pm) callAction(name string, payload interface{}) *string {
	if act, ok := pm2.actions.Load(name); ok {
		response := act.(*action).Callback(payload)
		return &response
	}
	return nil
}

// registerAction 将 action 注册到 pm2 中
func (pm2 *pm) registerAction(action *action) {
	msg := message{
		Type: "axm:action",
		Data: action,
	}
	pm2.tran.sendJson(msg)
}

// actionHandler 根据 pm2 触发的函数执行相应的 action
func (pm2 *pm) actionHandler(msg string, opts string) {
	defer func() {
		_ = recover()
	}()
	response := pm2.callAction(msg, opts)
	pm2.tran.sendJson(message{
		Type: "axm:reply",
		Data: map[string]interface{}{
			"action_name": msg,
			"return":      response,
		},
	})
}

// AddAction 注册 action 函数到 pm2 上
func (pm2 *pm) AddAction(name string, f actionFun) {
	if !pm2.isConnected() {
		return
	}
	action := &action{
		ActionName: name,
		ActionType: "custom",
		Callback:   f,
	}
	pm2.addAction(action)
	pm2.registerAction(action)
}
