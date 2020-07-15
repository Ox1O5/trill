package main

import (
	"fmt"
	"github.com/Ox1O5/trill/trill"
)

type pingRouter struct {
	trill.BaseRouter
}

func (this *pingRouter) Handle(request trill.IRequest)  {
	fmt.Println("Call pingRouter Handle")
	fmt.Println("receive from client : msgID = ", request.GetMsgID(),
		" data = ", string(request.GetData()))
	err := request.GetConnection().SendMsg(0, []byte("ping ping ping...\n"))
	if err != nil {
		fmt.Println("SendMsg error ", err)
	}
}

type helloRouter struct {
	trill.BaseRouter
}

func (h *helloRouter) Handle(request trill.IRequest) {
	fmt.Println("Call helloRouter Handle")
	fmt.Println("receive from client : msgID = ", request.GetMsgID(),
		" data = ", string(request.GetData()))

	err := request.GetConnection().SendMsg(1, []byte("hello hello hello...\n"))
	if err != nil {
		fmt.Println("SendMsg error ", err)
	}
}

func doConnectionBegin(conn trill.IConnection) {
	fmt.Println("DoConnecionBegin is Called ... ")
	fmt.Println("Set conn Name, Home done!")
	conn.SetProperty("Name", "0x105")
	conn.SetProperty("Home", "Ox105.github.io")
	err := conn.SendMsg(2, []byte("DoConnection BEGIN..."))
	if err != nil {
		fmt.Println(err)
	}
}

func doConnectionLost(conn trill.IConnection) {
	if name, err:= conn.GetProperty("Name"); err == nil {
		fmt.Println("Conn Property Name = ", name)
	}
	if home, err := conn.GetProperty("Home"); err == nil {
		fmt.Println("Conn Property Home = ", home)
	}
	fmt.Println("DoConnectionLost is Called ... ")
}


func main() {
	s := trill.NewServer("[trill v0.9]")
	s.SetOnConnStart(doConnectionBegin)
	s.SetOnConnStop(doConnectionLost)

	s.AddRouter(0, &pingRouter{})
	s.AddRouter(1, &helloRouter{})

	s.Server()
}
