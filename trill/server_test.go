package trill

import (
	"fmt"
	"io"
	"math/rand"
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
		pkt := NewPacket()
		msgID := uint32(rand.Intn(2))
		msg, err := pkt.Pack(NewMsgPacket(msgID, []byte("Trill v0.6 Client test message")))
		_, err = conn.Write(msg)
		if err != nil {
			fmt.Println("write error ", err)
			return
		}
		headData := make([]byte, pkt.GetHeadLen())
		_, err = io.ReadFull(conn, headData)
		if err != nil {
			fmt.Println("read head error")
			break
		}

		msgHead ,err := pkt.UnPack(headData)
		if err != nil {
			fmt.Println("server unpack err:", err)
			return
		}
		if msgHead.GetDataLen() > 0 {
			msg := msgHead.(*message)
			msg.data = make([]byte, msg.GetDataLen())
			_, err := io.ReadFull(conn, msg.data)
			if err != nil {
				fmt.Println("server unpack data error ", err)
				return
			}
			fmt.Println("==> Receive message : ID = ",
			msg.id, " len= ", msg.dataLen, " data= ", string(msg.data))
		}

		time.Sleep(time.Second)
	}
}

type pingRouter struct {
	baseRouter
}

//func (p *pingRouter) PreHandle(request IRequest) {
//	fmt.Println("Call Router PreHandle")
//	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping..."))
//	if err != nil {
//		fmt.Println("Call back ping error ", err)
//	}
//}

func (p *pingRouter) Handle(request IRequest) {
	fmt.Println("Call pingRouter Handle")
	fmt.Println("receive from client : msgID = ", request.GetMsgID(),
		" data = ", string(request.GetData()))
	err := request.GetConnection().SendMsg(0, []byte("ping ping ping...\n"))
	if err != nil {
		fmt.Println("SendMsg error ", err)
	}
}

//func (p *pingRouter) PostHandle(request IRequest) {
//	fmt.Println("Call Router PostHandle")
//	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping..."))
//	if err != nil {
//		fmt.Println("Call back ping error ", err)
//	}
//}

type helloRouter struct {
	baseRouter
}

func (h *helloRouter) Handle(request IRequest) {
	fmt.Println("Call helloRouter Handle")
	fmt.Println("receive from client : msgID = ", request.GetMsgID(),
		" data = ", string(request.GetData()))

	err := request.GetConnection().SendMsg(1, []byte("hello hello hello...\n"))
	if err != nil {
		fmt.Println("SendMsg error ", err)
	}
}

func TestServer(t *testing.T) {
	s := NewServer("[trill 0.6]")
	s.AddRouter(0, &pingRouter{})
	s.AddRouter(1, &helloRouter{})
	go ClientTest()
	s.Server()
}
