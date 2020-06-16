package Service

import (
	"fmt"
	"net"
	"strings"
	"time"
)

//定义基本结构体
type Service struct {
	conn        net.Conn
	recv        chan []byte
	sed         chan []byte
	writeFlag   chan bool
	er          chan bool
	serviceType int //0 本地服务 1 client 2 server 3 远程监听服务
}

func (server *Service) ping() {
	var pingStr = []byte("ping " + time.Now().String())
	_, _ = server.conn.Write(pingStr)
}

func (server *Service) pong() {
	var pongStr = []byte("pong " + time.Now().String())
	_, _ = server.conn.Write(pongStr)
}

func (server *Service) write() {
	for {
		var buf = make([]byte, 65535)
		select {
		case buf = <-server.sed:
			_, _ = server.conn.Write(buf)
		case <-server.writeFlag:
			//当读断开时 该获取模式不再写
			break
		}
	}
}

func (server *Service) read() {
	if server.serviceType == 1 {
		_ = server.conn.SetReadDeadline(time.Now().Add(time.Second * 20))
	} else if server.serviceType == 3 {
		_ = server.conn.SetReadDeadline(time.Now().Add(time.Second * 1))
	}
	var isHeart = false
	for {
		var buf = make([]byte, 65535)
		n, err := server.conn.Read(buf)
		if server.serviceType == 3 {
			_ = server.conn.SetReadDeadline(time.Time{})
		}
		if err != nil {
			if strings.Contains(err.Error(), "timeout") && !isHeart && server.serviceType == 1 {
				_ = server.conn.SetReadDeadline(time.Now().Add(time.Second * 3))
				fmt.Println(err.Error())
				server.ping()
				isHeart = true
				continue
			}
			server.writeFlag <- true
			server.er <- true
			break
		}
		if string(buf[:4]) == "pong" && server.serviceType == 1 {
			isHeart = false
			_ = server.conn.SetReadDeadline(time.Now().Add(time.Second * 20))
			continue
		} else if string(buf[:4]) == "ping" {
			server.pong()
			continue
		}
		server.recv <- buf[:n]
	}
}
