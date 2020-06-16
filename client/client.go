package client

import (
	"fmt"
	"go_nat_git/Service"
	"net"
)

func run() {
	var host = "127.0.0.1"
	//var host = "111.231.86.196"
	for {
		conn, _ := net.Dial("tcp", host+":10086")
		fmt.Println("dial conn success", conn.RemoteAddr().String())
		//recv := make(chan []byte)
		//sed := make(chan []byte)
		//run := make(chan bool)
		//er := make(chan bool, 1)     //错误管道 2端
		//writeFlag := make(chan bool) //写标志管道
		//clientType := 1
		var server = Service.Service{}
		//go read(conn, uR, false, er, writeFlag)
		//go write(conn, cR, false, writeFlag)
		//go handle(run, uR, cR, er, writeFlag)
		//<-run
	}
}
