package main

import (
	"fmt"
	"net"
	"time"
)

//基于本地转发 拓展的 nat 转发

func main() {
	conn, _ := net.Dial("tcp", "111.231.86.196:10086")
	fmt.Println("dial conn success")
	uR := make(chan []byte)
	cR := make(chan []byte)
	go read(conn, uR, false)
	go write(conn, cR)
	go ping(conn)
	for {
		run := make(chan bool)
		fmt.Println("recv data", time.Now().String())
		go handle(run, uR, cR)
		<-run
	}
}
func ping(conn net.Conn) {
	var pingStr = []byte("ping")
	//30s ping 一次
	var timeS = time.Second * 30
	for {
		_, _ = conn.Write(pingStr)
		time.Sleep(timeS)
	}
}
func handle(run chan bool, uR chan []byte, bR chan []byte) {
	var recv = make([]byte, 10240)
	fmt.Println("handle in")
	recv = <-uR
	fmt.Println("handle out")
	dial, _ := net.Dial("tcp", ":10084")
	fmt.Println("dial 127.0.1:10084")
	go write(dial, uR)
	go read(dial, bR, true)
	uR <- recv
	run <- true
}
func read(conn net.Conn, read chan []byte, isB bool) {
	if isB {
		//_ = conn.SetReadDeadline(time.Now().Add(time.Second * 10))
	}
	//_ = conn.SetReadDeadline(time.Now().Add(time.Second * 10))
	//var isHeart = false
	for {
		var buf = make([]byte, 10240)
		n, err := conn.Read(buf)
		fmt.Println("read", conn.LocalAddr().String())
		if err != nil {
			//if strings.Contains(err.Error(), "timeout") && !isHeart {
			//	_ = conn.SetReadDeadline(time.Now().Add(time.Second * 3))
			//	fmt.Println(err.Error())
			//	continue
			//}
			fmt.Println(err.Error())
			fmt.Println("close read")
			_ = conn.Close()
			break
		}
		if string(buf[:4]) == "pong" {
			continue
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
