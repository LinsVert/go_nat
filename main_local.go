package main

import (
	"fmt"
	"net"
	"time"
)

//基本的本地监听服务

func main() {
	listen, _ := net.Listen("tcp", ":10087")
	fmt.Println("start on 10087")
	for {
		conn, _ := listen.Accept()
		uR := make(chan []byte)
		bR := make(chan []byte)
		var run = make(chan bool)
		fmt.Println("connect", time.Now().String())
		go read(conn, uR)
		go write(conn, bR)
		go handle(run, uR, bR)
		<-run
		fmt.Println("close all")
	}

}
func handle(run chan bool, uR chan []byte, bR chan []byte) {
	var recv = make([]byte, 10240)
	fmt.Println("handle in")
	recv = <-uR
	fmt.Println("handle out")
	dial, _ := net.Dial("tcp", ":8000")
	fmt.Println("dial 127.0.1:8000")
	go write(dial, uR)
	go read(dial, bR)
	uR <- recv
	run <- true
}
func read(conn net.Conn, read chan []byte) {
	_ = conn.SetReadDeadline(time.Now().Add(time.Second * 10))
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

func write(conn net.Conn, write chan []byte) {
	for {
		var buf = make([]byte, 10240)
		select {
		case buf = <-write:
			fmt.Println("write Data", conn.LocalAddr())
			_ = conn.SetWriteDeadline(time.Now().Add(time.Second * 10))
			_, err := conn.Write(buf)
			if err != nil {
				fmt.Println("close write")
				_ = conn.Close()
				break
			}
		}
	}
}
