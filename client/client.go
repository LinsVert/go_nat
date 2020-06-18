package client

import (
	"fmt"
	"go_nat_git/Service"
	"net"
	"runtime"
)

func Run() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	var clientConfig = Service.GetClientConfig()
	for {
		conn, _ := net.Dial(clientConfig.RemoteConnect, net.JoinHostPort(clientConfig.RemoteAddress, clientConfig.RemotePort))
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
		go handle(server, run, clientConfig)
		<-run
	}
}

func handle(server Service.Service, run chan bool, clientConfig Service.ClientConfig) {
	var recv = make([]byte, 10240)
	recv = <-server.Recv
	run <- true
	dial, err := net.Dial(clientConfig.LocalConnect, net.JoinHostPort(clientConfig.LocalAddress, clientConfig.LocalPort))
	if err != nil {
		fmt.Println("Can't connect local server", err.Error())
		runtime.Goexit()
	}
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
