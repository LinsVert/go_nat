package main

import (
	"fmt"
	"net"
	"time"
)

//基于本地转发 拓展的 nat 转发

func main() {
	conn, _ := net.Dial("tcp", "127.0.0.1:10086")
	fmt.Println("dial conn success", conn.RemoteAddr().String())
	uR := make(chan []byte)
	cR := make(chan []byte)
	go read(conn, uR, false)
	go write(conn, cR, false)
	go ping(conn)
	for {
		//uR := make(chan []byte)
		//cR := make(chan []byte)
		//go read(conn, uR, false)
		//go write(conn, cR, false)
		run := make(chan bool)
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
	recv = <-uR
	fmt.Println("recv data", time.Now().String())
	dial, _ := net.Dial("tcp", ":8000")
	fmt.Println("dial", dial.RemoteAddr().String())
	go write(dial, uR, true)
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
		fmt.Println("read", conn.RemoteAddr().String(), n)
		if n < 100 && n > 0 {
			fmt.Println("Message", string(buf))
		}
		//if n == 0 {
		//	continue
		//}
		if err != nil {
			//if strings.Contains(err.Error(), "timeout") && !isHeart {
			//	_ = conn.SetReadDeadline(time.Now().Add(time.Second * 3))
			//	fmt.Println(err.Error())
			//	continue
			//}
			fmt.Println(err.Error())
			fmt.Println("close read", isB)
			_ = conn.Close()
			break
		}
		if string(buf[:4]) == "pong" {
			continue
		}
		read <- buf[:n]
	}
}

func write(conn net.Conn, write chan []byte, isB bool) {
	for {
		var buf = make([]byte, 10240)
		select {
		case buf = <-write:
			//_ = conn.SetWriteDeadline(time.Now().Add(time.Second * 10))
			fmt.Println("write Data", conn.RemoteAddr(), isB)
			_, err := conn.Write(buf)
			if err != nil {
				fmt.Println("close write")
				_ = conn.Close()
				break
			}
		}
	}
}
