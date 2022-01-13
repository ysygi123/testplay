package tcp

import (
	"sync"
)

type ConnectPoolManager struct {
	ConnectNum int
	Pool       sync.Pool
}

func NewConnectPoolManager() *ConnectPoolManager {
	cpm := new(ConnectPoolManager)
	cpm.Pool.New = func() interface{} {
		c := NewOneConnectHandler()
		c.Init()
		return c
	}
	return cpm
}

func (c *ConnectPoolManager) Producer(registerEvent *MyHandler) {
	x := (c.Pool.Get()).(*OneConnectHandler)
	x.MyChannel <- registerEvent
}

type OneConnectHandler struct {
	MyChannel chan *MyHandler
}

func NewOneConnectHandler() *OneConnectHandler {
	return &OneConnectHandler{
		MyChannel: make(chan *MyHandler, 1000),
	}
}

func (o *OneConnectHandler) Init() {
	go func() {
		for {
			x, ok := <-o.MyChannel
			if !ok {
				break
			}
			x.HandlerFunc(x.Params)
		}
	}()
}

type MyHandler struct {
	Params      interface{}
	HandlerFunc func(interface{}) interface{}
}
