package main

import (
	"bufio"
	"fmt"
	"net"
	"testplay/tcp"
)

func main() {
	tcp.Server()
}

func doSer(conn net.Conn) {
	for {
		bytes, err := bufio.NewReader(conn).ReadBytes('\n')
		if err != nil {
			fmt.Println("出现错误", err)
			return
		}
		fmt.Println(string(bytes))
	}
}
