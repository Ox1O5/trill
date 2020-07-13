package trill

import (
	"errors"
	"fmt"
	"net"
	"time"
)

type IServer interface {
	Start()
	Stop()
	Server()
}

type server struct {
	name      string
	ipVersion string
	ip        string
	port      int
}

func NewServer(name string) IServer {
	s := &server{
		name:      name,
		ipVersion: "tcp4",
		ip:        "0.0.0.0",
		port:      9090,
	}
	return s
}

func callBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
	fmt.Println("[connection handle] call back to client...")
	if _, err := conn.Write(data[:cnt]); err != nil {
		fmt.Println("write back error ", err)
		return errors.New("CallBackError")
	}
	return nil
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
			handleConn := NewConnection(conn, cid, callBackToClient)
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
