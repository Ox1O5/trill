package main

import (
	"fmt"
	"github.com/Ox1O5/trill/trill"
	"io"
	"net"
	"time"
)

func main() {
	fmt.Println("Client Test ... start")
	time.Sleep(3 * time.Second)
	conn, err := net.Dial("tcp", "127.0.0.1:9090")
	if err != nil {
		fmt.Println("client start error ", err)
		return
	}
	for {
		pkt := trill.NewPacket()
		//msgID := uint32(rand.Intn(2))
		msg, err := pkt.Pack(trill.NewMsgPacket(0, []byte("Trill v0.6 Client test message")))
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

		msgHead, err := pkt.UnPack(headData)
		if err != nil {
			fmt.Println("server unpack err:", err)
			return
		}
		if msgHead.GetDataLen() > 0 {
			msg := msgHead.(*trill.Message)
			msg.Data = make([]byte, msg.GetDataLen())
			_, err := io.ReadFull(conn, msg.Data)
			if err != nil {
				fmt.Println("server unpack data error ", err)
				return
			}
			fmt.Println("==> Receive message : ID = ",
				msg.ID, " len= ", msg.DataLen, " data= ", string(msg.Data))
		}

		time.Sleep(2*time.Second)

	}
}
