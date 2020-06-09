package logic

import (
	"fmt"
	"runtime"
)

func RunServer() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	var listenPort = "10086"
	var userPort = "10087"
	var client = Service{}
	var user = Service{}
	client.setListen("tcp", listenPort)
	client.listen()
	user.setListen("tcp", userPort)
	user.listen()
	fmt.Printf("监听成功 服务端口 %s, 业务端口 %s", listenPort, userPort)
	//开始监听服务
	//等待任意的连接
RE:
	go user.waitChanConn(user.listener)
	//阻塞等待客户端
	client.wait()
	fmt.Print("客户端连接成功")
	go client.receive()
	go client.send()
	for {
		//有两种情况
		//1.user 的数据传输
		//2 client 和 server的数据传输
		select {
			case conn := <- user.chanConn:
				user.conn = conn
				go user.receive()
				go user.send()
				go handle(client, user)
				goto RE
		}
	}
}