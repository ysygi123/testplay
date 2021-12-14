package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	l, err := net.Listen("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			return
		}

		go doSer(conn)
	}
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
