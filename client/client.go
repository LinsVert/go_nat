package client

import (
	"fmt"
	"go_nat_git/Service"
	"net"
)

func Run() {
	var host = "127.0.0.1"
	//var host = "111.231.86.196"
	for {
		conn, _ := net.Dial("tcp", host+":10086")
		fmt.Println("dial conn success", conn.RemoteAddr().String())
		recv := make(chan []byte)
		sed := make(chan []byte)
		run := make(chan bool)
		er := make(chan bool, 1)     //错误管道 2端
		writeFlag := make(chan bool) //写标志管道
		clientType := 1
		var server = Service.Service{Conn: conn, Recv: recv, Sed: sed, WriteFlag: writeFlag, Er: er, ServiceType: clientType}
		go server.Read()
		go server.Write()
		go handle(server, run)
		<-run
	}
}

func handle(server Service.Service, run chan bool) {
	var recv = make([]byte, 10240)
	recv = <-server.Recv
	run <- true
	dial, _ := net.Dial("tcp", ":8000")
	fmt.Println("dial conn success2", dial.RemoteAddr().String())
	recvL := make(chan []byte)
	sed := make(chan []byte)
	er := make(chan bool, 1)     //错误管道 2端
	writeFlag := make(chan bool) //写标志管道
	clientType := 0
	var local = Service.Service{Conn: dial, Recv: recvL, Sed: sed, WriteFlag: writeFlag, Er: er, ServiceType: clientType}
	go local.Read()
	go local.Write()
	local.Sed <- recv
	Service.Change(server, local)
}
