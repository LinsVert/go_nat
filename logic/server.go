package main

import (
	"fmt"
	"net"
	"runtime"
	"time"
)

func main() {
	listen, _ := net.Listen("tcp", ":10086")
	userL, _ := net.Listen("tcp", ":10087")
	fmt.Println("start on 10086")
	//go handleS(userL, uR,sR)
	//for {
	//	uR := make(chan []byte)
	//	sR := make(chan []byte)
	//	//go test(conn)
	//	go readS(conn, sR, false)
	//	go writeS(conn, uR)
	//	go handleS(userL, uR,sR)
	//	//for {
	//	//	uCon,_ := userL.Accept()
	//	//	go readS(uCon, uR)
	//	//	go writeS(uCon, sR)
	//	//}
	//}

	//var connChan = make(chan net.Conn)
	for {
		var connChan = make(chan net.Conn)
		//conn, _ := userL.Accept()
		go getConn(userL, connChan)
		uR := make(chan []byte)
		sR := make(chan []byte)
		er := make(chan bool, 1)     //错误管道 2端
		writeFlag := make(chan bool) //写标志管道
		conn1, _ := listen.Accept()
		fmt.Println("connect", time.Now().String())
		go readS(conn1, sR, false, er, writeFlag)
		go writeS(conn1, uR, false, writeFlag)
		//uRS := make(chan []byte)
		//sRS := make(chan []byte)
		var conn = <-connChan
		fmt.Println("get conn", conn)

		go readS(conn, uR, true, er, writeFlag)
		go writeS(conn, sR, true, writeFlag)
		//for {
		//	time.Sleep(time.Second * 5)
		//}
		go checkClose(conn1, conn, er)
		//go handleS(uR, sR, sRS, uRS)
	}
}
func checkClose(conn1 net.Conn, conn2 net.Conn, er chan bool) {
	for {
		select {
		case <-er:
			_ = conn1.Close()
			_ = conn2.Close()
			runtime.Goexit()
		}
	}
}

func getConn(listener net.Listener, connChan chan net.Conn) {
	conn, _ := listener.Accept()
	connChan <- conn
}

func handleS(uR chan []byte, sR chan []byte, sRS chan []byte, uRS chan []byte) {
	for {
		var buf = make([]byte, 10240)
		select {
		case buf = <-sR:
			sRS <- buf
		case buf = <-uRS:
			uR <- buf
		}
	}
}

func readS(conn net.Conn, read chan []byte, isU bool, er chan bool, writeFlag chan bool) {
	for {
		if isU {
			_ = conn.SetReadDeadline(time.Now().Add(time.Second * 1))
		}
		var buf = make([]byte, 10240)
		n, err := conn.Read(buf)
		fmt.Println("read", conn.LocalAddr().String())
		if n < 100 && n > 0 {
			fmt.Println("Message", string(buf))
		}
		if err != nil {
			fmt.Println("close read", isU)
			//_ = conn.Close()
			er <- true
			writeFlag <- true
			break
		}
		if isU {
			_ = conn.SetReadDeadline(time.Time{})
		}
		if string(buf[:4]) == "ping" {
			_, _ = conn.Write([]byte("pong " + time.Now().String()))
			continue
		}
		fmt.Println(n, "read at", isU)
		read <- buf[:n]
	}
}

func writeS(conn net.Conn, read chan []byte, isU bool, writeFlag chan bool) {
	for {
		var buf = make([]byte, 10240)
		//if isU {
		//	_ = conn.SetWriteDeadline(time.Now().Add(time.Second * 10))
		//}
		select {
		case buf = <-read:
			n, err := conn.Write(buf)
			if n < 100 && n > 0 {
				fmt.Println("Message write", string(buf))
			}
			fmt.Println(n, "write at", isU)
			if err != nil {
				fmt.Println("close write2")
				_ = conn.Close()
				break
			}
			//if isU {
			//	_ = conn.SetWriteDeadline(time.Now())
			//}
			if isU {
				//fmt.Println("Message write to U", string(buf))
			}
		case <-writeFlag:
			break
		}

	}
}
