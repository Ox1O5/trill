package trill

import (
	"net",
	"fmt"
)

type IServer struct {
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
		port : "9090",
	}
	return s
}

func (s *server) Start() {
	fmt.Printf("[start] server listenner at ip: %s : %d , is starting\n", s.ip, s.port)
	go func() {
		add ,err := net.ResolveTCPAddr(s.ipVersion, fmt.Sprintf("%s:%d", s.ip, s.port))
		if err != nil {
			fmt.Println("resolve tcp addr err: ", err)
			return
		}
		listenner, err := net.Lis
	}
}