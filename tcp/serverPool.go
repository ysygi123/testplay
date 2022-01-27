package tcp

import (
	"net"
	"sync"
	"testplay/utils"
)

type ConnectPoolManager struct {
	ConnectNum int
	Pool       sync.Pool
	//handle func(c net.Conn)
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

func (c *ConnectPoolManager) Producer(conn net.Conn) {
	x := (c.Pool.Get()).(*OneConnectHandler)
	x.MyChannel <- conn
}

type OneConnectHandler struct {
	MyChannel chan net.Conn
	doubleReadWritePool sync.Pool
}

func NewOneConnectHandler() *OneConnectHandler {
	och := new(OneConnectHandler)
	och.doubleReadWritePool.New = func() interface{} {
		c := NewReadWrite()
		c.ReadInit()
		c.WriteInit()
		return c
	}
	och.MyChannel = make(chan net.Conn, 1000)
	return och
}

func (o *OneConnectHandler) Init() {
	go func() {
		for {
			conn, ok := <-o.MyChannel
			if !ok {
				break
			}
			o.MyServer(conn)
		}
	}()
}

func (o *OneConnectHandler)MyServer(conn net.Conn) {
	defer conn.Close()
	//声明一个临时缓冲区，用来存储被截断的数据
	tmpBuffer := make([]byte, 0)
	//声明一个管道用于接收解包的数据
	readerChannel := make(chan []byte, 16)
	go reader(readerChannel)

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			return
		}

		tmpBuffer = utils.UnPack(append(tmpBuffer, buffer[:n]...), readerChannel)
	}
}


type DoubleReadWrite struct {
	ReadChannel  chan net.Conn
	WriteChannel chan net.Conn
}

func NewReadWrite() *DoubleReadWrite {
	return &DoubleReadWrite{
		ReadChannel:  make(chan net.Conn, 1000),
		WriteChannel: make(chan net.Conn, 1000),
	}
}

func (d *DoubleReadWrite)ReadInit()  {
	go func() {
		for {
			conn, ok := <-d.ReadChannel
			if !ok {
				break
			}
			d.Read(conn)
		}
	}()
}

func (d *DoubleReadWrite)Read(conn net.Conn)  {
	for {
		select {
		//case data := <-readerChannel:
		//	fmt.Println(utils.B2S(data))
		}
	}
}

func (d *DoubleReadWrite)WriteInit()  {
	go func() {
		for {
			conn, ok := <-d.WriteChannel
			if !ok {
				break
			}
			d.Write(conn)
		}
	}()
}

func (d *DoubleReadWrite)Write(conn net.Conn)  {

}