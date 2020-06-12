package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	listen, _ := net.Listen("tcp", ":10086")
	fmt.Println("start on 10086")
	for {
		conn, _ := listen.Accept()
		uR := make(chan []byte)
		fmt.Println("connect", time.Now().String())
		go test(conn)
		go readS(conn, uR)
		go writeS(conn, uR)
	}
}

func test(conn net.Conn) {
	for {
		time.Sleep(time.Second * 10)
		_, _ = conn.Write([]byte("hello world on 10086"))
	}
}

func readS(conn net.Conn, read chan []byte) {
	//_ = conn.SetReadDeadline(time.Now().Add(time.Second * 10))
	for {
		var buf = make([]byte, 10240)
		n, err := conn.Read(buf)
		fmt.Println("read", conn.LocalAddr().String())
		if err != nil {
			fmt.Println("close read")
			_ = conn.Close()
			break
		}

		read <- buf[:n]
	}
}

func writeS(conn net.Conn, read chan []byte) {
	for {
		var buf = make([]byte, 10240)
		select {
		case buf = <-read:
			if string(buf) == "ping" {
				_, err := conn.Write([]byte("pong " + time.Now().String()))
				if err != nil {
					fmt.Println("close write")
					_ = conn.Close()
					break
				}
			} else {
				fmt.Println(string(buf))
			}

		}
	}
}
