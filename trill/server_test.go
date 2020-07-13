package trill

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func ClientTest() {
	fmt.Println("Client Test ... start")
	time.Sleep(3 * time.Second)
	conn, err := net.Dial("tcp", "127.0.0.1:9090")
	if err != nil {
		fmt.Println("client start error ", err)
		return
	}
	for {
		_, err := conn.Write([]byte("Hello world"))
		if err != nil {
			fmt.Println("write error ", err)
			return
		}
		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read from buffer error ", err)
			return
		}
		fmt.Printf("server call back : %s\n", buf[:cnt])
		time.Sleep(time.Second)
	}
}

type pingRouter struct {
	baseRouter
}

func (p *pingRouter) PreHandle(request IRequest) {
	fmt.Println("Call Router PreHandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping..."))
	if err != nil {
		fmt.Println("Call back ping error ", err)
	}
}

func (p *pingRouter) Handle(request IRequest) {
	fmt.Println("Call Router Handle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping ping ping..."))
	if err != nil {
		fmt.Println("Call back ping error ", err)
	}
}

func (p *pingRouter) PostHandle(request IRequest) {
	fmt.Println("Call Router PostHandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping..."))
	if err != nil {
		fmt.Println("Call back ping error ", err)
	}
}

func TestServer(t *testing.T) {
	s := NewServer("[trill 0.3]")
	s.AddRouter(&pingRouter{})
	go ClientTest()

	s.Server()
}
