package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	b := []byte("这里是一段话，用来标志看看有没有很奇怪的东西会编出来 []\r\n")
	conn, err := net.Dial("tcp", "127.0.0.1:9999")
	if err != nil {
		panic(err)
	}
	t := time.Now()
	for i := 0; i < 1500000; i++ {
		conn.Write(b)
	}
	fmt.Println(time.Since(t))

	//t := time.Now()
	//wg := sync.WaitGroup{}
	//wg.Add(1)
	//
	//go func() {
	//	for i := 0; i < 500000; i++ {
	//		conn.Write(b)
	//	}
	//	wg.Done()
	//}()
	//wg.Add(1)
	//go func() {
	//	for i := 0; i < 500000; i++ {
	//		conn.Write(b)
	//	}
	//	wg.Done()
	//}()
	//wg.Add(1)
	//go func() {
	//	for i := 0; i < 500000; i++ {
	//		conn.Write(b)
	//	}
	//	wg.Done()
	//}()
	//wg.Wait()
	//fmt.Println(time.Since(t))
}
