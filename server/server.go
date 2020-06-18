package server

import (
	"fmt"
	"go_nat_git/Service"
	"net"
	"runtime"
	"time"
)

func Run() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	var serverConfig = Service.GetServerConfig()
	listen, _ := net.Listen(serverConfig.ServiceListenConnect, ":"+serverConfig.ServerListenPort)
	userL, _ := net.Listen(serverConfig.UserListenConnect, ":"+serverConfig.UserListenPort)
	fmt.Println("start on", serverConfig.ServerListenPort, serverConfig.UserListenPort)
	for {
		var connChan = make(chan net.Conn)
		go Service.GetConn(userL, connChan)
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
