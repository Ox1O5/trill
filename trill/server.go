package trill

import (
	"fmt"
	"github.com/Ox1O5/trill/utils"
	"net"
	"time"
)

type IServer interface {
	Start()
	Stop()
	Server()
	AddRouter(router IRouter)
}

type server struct {
	name      string
	ipVersion string
	ip        string
	port      int
	router    IRouter
}

func NewServer(name string) IServer {
	utils.GlobalObject.Load()
	s := &server{
		name:      utils.GlobalObject.Name,
		ipVersion: "tcp4",
		ip:        utils.GlobalObject.Host,
		port:      utils.GlobalObject.TcpPort,
		router:    nil,
	}
	return s
}


func (s *server) Start() {
	fmt.Printf("[start] server listenner at ip: %s : %d , is starting\n", s.ip, s.port)
	go func() {
		addr, err := net.ResolveTCPAddr(s.ipVersion, fmt.Sprintf("%s:%d", s.ip, s.port))
		if err != nil {
			fmt.Println("resolve tcp addr err: ", err)
			return
		}
		listener, err := net.ListenTCP(s.ipVersion, addr)
		if err != nil {
			fmt.Println("listen", addr, "error", err)
			return
		}
		fmt.Println("start trill ", s.name, " success, listening at ", addr)
		var cid uint32
		cid = 0
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept error ", err)
				continue
			}
			handleConn := NewConnection(conn, cid, s.router)
			cid++
			go handleConn.Start()
		}
	}()
}

func (s *server) Stop() {
	fmt.Println("[stop] server name ", s.name)
}

func (s *server) Server() {
	s.Start()
	for {
		time.Sleep(10 * time.Second)
	}
}

func (s *server) AddRouter(router IRouter) {
	s.router = router
	fmt.Println("Add router success")
}
