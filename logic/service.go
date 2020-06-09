package logic

import (
	"fmt"
	"net"
	"os"
)

type Service struct {
	conn net.Conn
	chanConn chan net.Conn
	dialConn net.Conn
	listener net.Listener
	listenerPort string
	listenerType string
	rec chan []byte
	sed chan []byte
	err error
	to int
	localAddr string
	localPort string
	remoteAddr string
	remotePort string
	localConnType string
	remoteConnType string
}


func (service *Service) setLocalAddr(localAddr string, localPort string, localConnType string) {
	service.localAddr = localAddr
	service.localPort = localPort
	service.localConnType = localConnType
}

func (service *Service) setRemoteAddr(remoteAddr string, remotePort string, remoteConnType string) {
	service.remoteAddr = remoteAddr
	service.remotePort = remotePort
	service.remoteConnType = remoteConnType
}

func (service * Service) setListen(listenerType string, listenerPort string) {
	service.listenerType = listenerType
	service.listenerPort = ":" + listenerPort
	fmt.Printf(service.listenerPort)
}
func (service *Service) dial(to int) {

	if to == 0 {
		service.dialConn, service.err = net.Dial(service.localConnType, net.JoinHostPort(service.localAddr, service.localPort))
	} else if to == 1 {
		service.dialConn, service.err = net.Dial(service.remoteConnType, net.JoinHostPort(service.remoteAddr, service.remotePort))
	} else {
		service.dialConn, service.err = net.Dial(service.localConnType, net.JoinHostPort(service.localAddr, service.localPort))
	}
	service.errMsg(service.err)
}

func (service *Service) listen() {
	service.listener, service.err = net.Listen(service.listenerType, service.listenerPort)
	service.errMsg(service.err)
}

func (service *Service) errMsg (err error){
	if err != nil {
		fmt.Print(err, "\n")
		os.Exit(0)
	}
}

//发送
func (service *Service) send() {
	for {
		var send []byte = make([]byte, 10240)
		select {
			case send = <- service.sed:
				_, err := service.conn.Write(send)
				service.errMsg(err)
				if err != nil {
					break
				}
		}
	}
}
//接收
func (service *Service) receive() {
	for {
		var rec = make([]byte, 10240)
		n, err := service.conn.Read(rec)
		service.errMsg(err)
		if err != nil {
			//中断 todo
			break;
		}
		service.initData()
		service.rec <- rec[:n]
	}
}

func (service *Service) wait() net.Conn {
	//建立连接
	service.conn, service.err = service.listener.Accept()
	service.errMsg(service.err)
	fmt.Print("已经建立连接")
	return service.conn
}
func (service *Service) waitChanConn(listener net.Listener) {
	Conn , err := listener.Accept()
	service.errMsg(err)
	service.initConnChan()
	service.chanConn <- Conn
}

func (server *Service) initConnChan() chan net.Conn {
	if server.chanConn != nil {
		return server.chanConn
	}
	server.chanConn = make(chan net.Conn)
	return server.chanConn
}
func (service *Service) initData() {
	if service.rec == nil {
		service.rec = make(chan []byte)
	}
	if service.sed == nil {
		service.sed = make(chan []byte)
	}
}


func (service *Service) heartbeat() {
	//todo
	var heart []byte = make([]byte, 10240)
	heart[0] = 1
	service.send()
}

//数据交换
func handle(client Service, user Service) {
	for {
		select {
			case rec := <- client.rec:
				user.sed <- rec
			case rec := <- user.rec:
				client.sed <-rec
		}
	}
}