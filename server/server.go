package server

import (
	"fmt"
	"go_nat_git/Service"
	"net"
	"time"
)

func Run() {
	listen, _ := net.Listen("tcp", ":10086")
	userL, _ := net.Listen("tcp", ":10087")
	fmt.Println("start on 10086")
	for {
		var connChan = make(chan net.Conn)
		go getConn(userL, connChan)
		recv := make(chan []byte)
		sed := make(chan []byte)
		er := make(chan bool, 1)     //错误管道 2端
		writeFlag := make(chan bool) //写标志管道
		clientType := 2
		clientConn, _ := listen.Accept() //等待客户端连接
		fmt.Println("connect in ", time.Now().String(), clientConn.LocalAddr().String())
		client := Service.Service{Conn: clientConn, Recv: recv, Sed: sed, WriteFlag: writeFlag, Er: er, ServiceType: clientType}
		go client.Read()
		go client.Write()
		var conn = <-connChan
		fmt.Println("user in ", time.Now().String(), conn.RemoteAddr().String())
		recv = make(chan []byte)
		sed = make(chan []byte)
		er = make(chan bool, 1)     //错误管道 2端
		writeFlag = make(chan bool) //写标志管道
		clientType = 3
		user := Service.Service{Conn: conn, Recv: recv, Sed: sed, WriteFlag: writeFlag, Er: er, ServiceType: clientType}
		go user.Read()
		go user.Write()
		go Service.Change(client, user)
	}
}

func getConn(listener net.Listener, connChan chan net.Conn) {
	conn, _ := listener.Accept()
	connChan <- conn
}
