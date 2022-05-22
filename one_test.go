package main

import (
	"bytes"
	"sync"
	"testing"
	"testplay/utils"
)

type TestMy777 struct {
	Name string
}

//func Test_A1(t *testing.T) {
//	cache := make(map[int]**TestMy777)
//	tt := &TestMy777{
//		Name: "",
//	}
//	cache[0] = &tt
//	x := new(sync.Mutex)
//	cache2 := make(map[int]*TestMy777)
//	go func() {
//		time.Sleep(1 * time.Second)
//		x.Lock()
//		delete(cache, 0)
//		x.Unlock()
//	}()
//	cache2[0] = *cache[0]
//	ticker := time.NewTicker(200 * time.Millisecond)
//	for {
//		select {
//		case <-ticker.C:
//			x.Lock()
//			_, ok := cache2[0]
//			x.Unlock()
//
//			if ok {
//				fmt.Println("cache2 还在")
//			} else {
//				fmt.Println("cache2 不在了")
//				goto END
//			}
//
//			x.Lock()
//			_, ok = cache[0]
//			x.Unlock()
//
//			if ok {
//				fmt.Println("cache1 还在")
//			} else {
//				fmt.Println("cache1 不在了")
//				//goto END
//			}
//		}
//	}
//END:
//	fmt.Println("over")
//}

func Benchmark_Pool(b *testing.B) {
	p := sync.Pool{
		New: func() interface{} {
			return &bytes.Buffer{}
		},
	}
	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 10000000; j++ {
				x := p.Get().(*bytes.Buffer)
				x.Write(utils.S2B("123"))
				x.Reset()
				p.Put(x)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func Benchmark_NotPool(b *testing.B) {
	_ = sync.Pool{
		New: func() interface{} {
			return &bytes.Buffer{}
		},
	}
	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 10000000; j++ {
				x := &bytes.Buffer{}
				x.Write(utils.S2B("123"))
				x.Reset()
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
