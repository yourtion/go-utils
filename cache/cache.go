// cache 是缓存系统的整体抽象
package cache

import (
	"sync/atomic"
)

// 返回 Cache 的统计信息
type Status struct {
	Gets        uint64 `json:"gets"`
	Hits        uint64 `json:"hits"`
	MaxItemSize int    `json:"max"`
	CurrentSize int    `json:"current"`
}

// this is a interface which defines some common functions
type Cache interface {
	Set(key string, value interface{})
	Get(key string) (interface{}, bool)
	Delete(key string)
	Status() *Status
}

// An atomicUint is an int64 to be accessed atomically.
type atomicUint uint64

func (i *atomicUint) add(n uint64) {
	atomic.AddUint64((*uint64)(i), n)
}

func (i *atomicUint) get() uint64 {
	return atomic.LoadUint64((*uint64)(i))
}
