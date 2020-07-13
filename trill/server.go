package trill

import (
	"net"
	"fmt"
	"time"
)

type IServer interface {
	Start()
	Stop()
	Server()
}

type server struct{
	name string
	ipVersion string
	ip string
	port int
}

func NewServer(name string)  IServer{
	s := &server {
		name : name,
		ipVersion : "tcp4",
		ip : "0.0.0.0",
		port : 9090,
	}
	return s
}

func (s *server) Start() {
	fmt.Printf("[start] server listenner at ip: %s : %d , is starting\n", s.ip, s.port)
	go func() {
		addr ,err := net.ResolveTCPAddr(s.ipVersion, fmt.Sprintf("%s:%d", s.ip, s.port))
		if err != nil {
			fmt.Println("resolve tcp addr err: ", err)
			return
		}
		listenner, err := net.ListenTCP(s.ipVersion, addr )
		if err != nil {
			fmt.Println("listen", addr ,"error",err)
			return
		}
		fmt.Println("start trill ", s.name, " success, listenning at ", addr)
		for {
			conn, err := listenner.Accept()
			if err != nil {
				fmt.Println("Accept error ", err)
				continue
			}
			go func() {
				for {
					buf := make([]byte, 512)
					cnt, err := conn.Read(buf)
					if err != nil {
						fmt.Println("receive buffer error ", err)
						continue
					}
					if _, err := conn.Write(buf[:cnt]); err != nil {
						fmt.Println("Write back error ", err)
						continue
					}
				}
			}()
		}
	}()
}

func (s *server) Stop()  {
	fmt.Println("[stop] server name ", s.name)
}

func (s *server) Server()  {
	s.Start()
	for {
		time.Sleep(10*time.Second)
	}
	
}