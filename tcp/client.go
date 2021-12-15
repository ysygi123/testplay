package tcp

import (
	"fmt"
	"net"
	"os"
	"testplay/utils"
	"time"
)

func sender(conn net.Conn) {
	for i := 0; i < 1000; i++ {
		message := utils.S2B("我可娶你娘亲")
		conn.Write(utils.Pack(message))
	}
}

func Client() {
	server := "127.0.0.1:9988"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	defer conn.Close()
	fmt.Println("connect success")
	go sender(conn)
	for {
		time.Sleep(1 * 1e9)
	}
}
