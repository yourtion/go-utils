package memo

import (
	"sync"
)

// 数据获取方法
type Func func(key string) (interface{}, error)

// 数据获取结果
type result struct {
	value interface{}
	err   error
}

type entry struct {
	res   result
	ready chan struct{} // closed when res is ready
}

/*
1. 获取互斥锁来保护共享变量cache map，查询map中是否存在指定key的value
2. 如果不存在那么分配空间插入新值，释放互斥锁。
3. 如果存在且其值没有写入完即其他goroutine在调用f这个慢函数，
4. goroutine必须等待值ready之后才能读到key的value
5. 如果没有key对应的value，需要向map中插入一个没有ready的entry，
6. 当前正在调用的goroutine就需要负责调用慢函数更新result以及向其他goroutine
7. 广播（关闭ready）result已经可读了
*/
type Memo struct {
	f     Func
	mu    sync.Mutex
	cache map[string]*entry
}

func NewMemo(f Func) *Memo {
	return &Memo{f: f, cache: make(map[string]*entry)}
}

func (mc *Memo) Get(key string) (interface{}, error) {
	mc.mu.Lock()
	e := mc.cache[key]
	if e == nil { // e为空
		// This is the first request for this key
		e = &entry{ready: make(chan struct{})}
		mc.cache[key] = e
		mc.mu.Unlock()

		e.res.value, e.res.err = mc.f(key) // 执行耗时函数
		close(e.ready)                     // broadcast ready condition
		// 当执行结束后删除缓存key
		mc.mu.Lock()
		mc.cache[key] = nil
		mc.mu.Unlock()
	} else {
		mc.mu.Unlock()
		<-e.ready // 阻塞，直到ready关闭
	}
	return e.res.value, e.res.err
}
