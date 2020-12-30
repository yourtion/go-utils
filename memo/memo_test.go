package memo

import (
	"math/rand"
	"sync"
	"testing"
	"time"
)

func getData(key string) (interface{}, error) {
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(100)
	time.Sleep(time.Duration(i) * time.Millisecond)
	println(key, i)
	return i, nil
}

func TestMemo(t *testing.T) {
	memo := NewMemo(getData)
	result := make([]int, 100)
	wg := sync.WaitGroup{}
	wg.Add(len(result))
	// 多个进程一并获取结果应当一致
	for i := range result {
		go func(k int) {
			ret, _ := memo.Get("i")
			result[k] = ret.(int)
			wg.Done()
		}(i)
	}
	wg.Wait()
	r := result[0]
	for i := range result {
		if result[i] != r {
			t.Fatalf("result[%d] != %d", i, r)
		}
	}
}
