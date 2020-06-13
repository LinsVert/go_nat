package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	listen, _ := net.Listen("tcp", ":10086")
	userL, _ := net.Listen("tcp", ":10087")
	fmt.Println("start on 10086")
	conn1, _ := listen.Accept()
	fmt.Println("connect", time.Now().String())
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
	uR := make(chan []byte)
	sR := make(chan []byte)
	go readS(conn1, sR, false)
	go writeS(conn1, uR, false)
	var connChan = make(chan net.Conn)
	for {
		//conn, _ := userL.Accept()
		go getConn(userL, connChan)
		//uRS := make(chan []byte)
		//sRS := make(chan []byte)
		var conn = <-connChan
		fmt.Println("get conn", conn)

		go readS(conn, uR, true)
		go writeS(conn, sR, true)
		//for {
		//	time.Sleep(time.Second * 5)
		//}

		//go handleS(uR, sR, sRS, uRS)
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

func readS(conn net.Conn, read chan []byte, isU bool) {
	for {
		if isU {
			//_ = conn.SetReadDeadline(time.Now().Add(time.Second * 10))
		}
		var buf = make([]byte, 10240)
		n, err := conn.Read(buf)
		fmt.Println("read", conn.LocalAddr().String())
		if n < 100 && n > 0 {
			fmt.Println("Message", string(buf))
		}
		if err != nil {
			fmt.Println("close read", isU)
			_ = conn.Close()
			break
		}
		if string(buf[:4]) == "ping" {
			_, _ = conn.Write([]byte("pong " + time.Now().String()))
			continue
		}
		fmt.Println(n, "read at", isU)
		read <- buf[:n]
	}
}

func writeS(conn net.Conn, read chan []byte, isU bool) {
	for {
		var buf = make([]byte, 10240)
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
		}

	}
}
