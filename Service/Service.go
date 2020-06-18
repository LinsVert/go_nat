package Service

import (
	"fmt"
	"net"
	"runtime"
	"strings"
	"time"
)

//定义基本结构体
type Service struct {
	Conn        net.Conn
	Recv        chan []byte
	Sed         chan []byte
	WriteFlag   chan bool
	Er          chan bool
	ServiceType int //0 本地服务 1 client 2 server 3 远程监听服务
}

func (server *Service) Ping() {
	var pingStr = []byte("Ping " + time.Now().String())
	_, _ = server.Conn.Write(pingStr)
}

func (server *Service) Pong() {
	var pongStr = []byte("Pong " + time.Now().String())
	_, _ = server.Conn.Write(pongStr)
}

func (server *Service) Write() {
	for {
		var buf = make([]byte, 65535)
		select {
		case buf = <-server.Sed:
			lens := len(buf)
			n, err := server.Conn.Write(buf)
			if server.ServiceType == 1 && lens != n {
				fmt.Println("write data", n, "get len", lens)
				if err != nil {
					fmt.Println(err.Error(), "write to server error")
				}
			}
		case <-server.WriteFlag:
			//当读断开时 该获取模式不再写
			break
		}
	}
}

func (server *Service) Read() {
	if server.ServiceType == 1 {
		_ = server.Conn.SetReadDeadline(time.Now().Add(time.Second * 20))
	} else if server.ServiceType == 3 {
		_ = server.Conn.SetReadDeadline(time.Now().Add(time.Second * 1))
	}
	var isHeart = false
	var finishRead = 0
	var limit = 100
	for {
		var buf = make([]byte, 65535)
		fmt.Println("wait data in ", server.ServiceType, server.Conn.LocalAddr().String(), server.Conn.RemoteAddr().String())
		n, err := server.Conn.Read(buf)
		if server.ServiceType == 0 {
			fmt.Println("recv data on ", n, server.ServiceType, server.Conn.LocalAddr().String(), server.Conn.RemoteAddr().String())
		}
		//fmt.Println("recv data on ", n, server.ServiceType)
		if server.ServiceType == 3 {
			_ = server.Conn.SetReadDeadline(time.Time{})
		}
		if err != nil {
			if strings.Contains(err.Error(), "timeout") && !isHeart && server.ServiceType == 1 {
				_ = server.Conn.SetReadDeadline(time.Now().Add(time.Second * 3))
				fmt.Println(err.Error())
				server.Ping()
				isHeart = true
				continue
			}
			if server.ServiceType == 1 {
				server.Recv <- []byte("***")
			}
			if server.ServiceType == 0 && err.Error() == "EOF" && finishRead < limit {
				//当数据请求一直空是 循序判空n次
				fmt.Println(err.Error(), "test1", finishRead, server.Conn.LocalAddr().String(), server.Conn.RemoteAddr().String())
				finishRead++
				//continue
			}
			server.WriteFlag <- true
			server.Er <- true
			break
		}
		if string(buf[:4]) == "Pong" && server.ServiceType == 1 {
			isHeart = false
			_ = server.Conn.SetReadDeadline(time.Now().Add(time.Second * 20))
			continue
		} else if string(buf[:4]) == "Ping" {
			server.Pong()
			continue
		}
		server.Recv <- buf[:n]
	}
}

func Change(service Service, next Service) {
	for {
		var buf = make([]byte, 65535)
		select {
		case buf = <-service.Recv:
			if string(buf[:3]) == "***" {
				continue
			}
			fmt.Println("Change Data in Service S", service.ServiceType, service.Conn.LocalAddr().String(), service.Conn.RemoteAddr().String())
			next.Sed <- buf
		case buf = <-next.Recv:
			fmt.Println("Change Data in Next N", next.ServiceType, next.Conn.LocalAddr().String(), next.Conn.RemoteAddr().String())
			service.Sed <- buf
		case <-service.Er:
			fmt.Println("close in this serviceSSSSSS:", service.ServiceType, service.Conn.LocalAddr().String(), service.Conn.RemoteAddr().String())
			_ = service.Conn.Close()
			_ = next.Conn.Close()
			runtime.Goexit()
		case <-next.Er:
			fmt.Println("close in this nextNNNNN:", next.ServiceType, next.Conn.LocalAddr().String(), next.Conn.RemoteAddr().String())
			_ = service.Conn.Close()
			_ = next.Conn.Close()
			runtime.Goexit()
		}
	}
}

func GetConn(listener net.Listener, connChan chan net.Conn) {
	conn, _ := listener.Accept()
	connChan <- conn
}
