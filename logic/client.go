package logic

import "fmt"

//client logic

func RunClient() {
	var local = "127.0.0.1"
	var localPort = "8000"
	//var listenPort = "10087"
	var remotePort = "10086"
	var server = Service{}
	var localServer = Service{}
	localServer.setLocalAddr(local, localPort, "TCP")
	server.setLocalAddr(local, localPort, "TCP")
	server.setRemoteAddr(local, remotePort, "TCP")
	server.dial(1)
	server.conn = server.dialConn
	fmt.Print("连接远端服务器成功")
	for {
		go server.receive()
		go server.send()
		localServer.dial(0)
		localServer.conn = localServer.dialConn
		go handle(server, server)
	}
}