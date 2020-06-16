package main

import (
	"fmt"
	"net"
	"runtime"
	"strings"
	"time"
)

//基于本地转发 拓展的 nat 转发

func main() {

	var host = "127.0.0.1"
	//var host = "111.231.86.196"
	for {
		conn, _ := net.Dial("tcp", host+":10086")
		fmt.Println("dial conn success", conn.RemoteAddr().String())
		uR := make(chan []byte)
		cR := make(chan []byte)
		run := make(chan bool)
		er := make(chan bool, 1)     //错误管道 2端
		writeFlag := make(chan bool) //写标志管道
		go read(conn, uR, false, er, writeFlag)
		go write(conn, cR, false, writeFlag)
		go handle(run, uR, cR, er, writeFlag)
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
func pingOnce(conn net.Conn) {
	var pingStr = []byte("ping")
	_, _ = conn.Write(pingStr)
}
func handle(run chan bool, uR chan []byte, bR chan []byte, er chan bool, writeFlag chan bool) {
	var recv = make([]byte, 10240)
	recv = <-uR
	run <- true
	fmt.Println("recv data", time.Now().String())
	dial, _ := net.Dial("tcp", ":8000")
	fmt.Println("dial", dial.RemoteAddr().String())
	go write(dial, uR, true, writeFlag)
	go read(dial, bR, true, er, writeFlag)
	uR <- recv
	for {
		select {
		case <-er:
			_ = dial.Close()
			runtime.Goexit()

		}
	}
}
func read(conn net.Conn, read chan []byte, isB bool, er chan bool, writeFlag chan bool) {
	if !isB {
		_ = conn.SetReadDeadline(time.Now().Add(time.Second * 20))
	}
	//_ = conn.SetReadDeadline(time.Now().Add(time.Second * 10))
	var isHeart = false
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
			if strings.Contains(err.Error(), "timeout") && !isHeart && !isB {
				_ = conn.SetReadDeadline(time.Now().Add(time.Second * 3))
				fmt.Println(err.Error())
				pingOnce(conn)
				isHeart = true
				continue
			}
			fmt.Println(err.Error())
			fmt.Println("close read", isB)
			//_ = conn.Close()
			if !isB {
				read <- []byte("0")
			}
			er <- true
			writeFlag <- true
			break
		}
		if string(buf[:4]) == "pong" && !isB {
			isHeart = false
			_ = conn.SetReadDeadline(time.Now().Add(time.Second * 20))
			continue
		}
		if !isB {
			//_ = conn.SetReadDeadline(time.Now())
		}
		read <- buf[:n]
	}
}

func write(conn net.Conn, write chan []byte, isB bool, writeFlag chan bool) {
	for {
		//if isB {
		//	_ = conn.SetWriteDeadline(time.Now().Add(time.Second * 10))
		//}
		var buf = make([]byte, 10240)
		select {
		case buf = <-write:
			//_ = conn.SetWriteDeadline(time.Now().Add(time.Second * 10))
			fmt.Println("write Data", conn.RemoteAddr(), isB)
			if isB {
				//fmt.Println(string(buf))
			}
			_, _ = conn.Write(buf)
		//if err != nil {
		//	fmt.Println("close write")
		//	//_ = conn.Close()
		//	break
		//}
		//if isB {
		//	_ = conn.SetWriteDeadline(time.Now())
		//}
		case <-writeFlag:
			fmt.Println("write flag in ")
			break
		}
	}
}
