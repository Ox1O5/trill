package trill

import (
	"errors"
	"fmt"
	"io"
	"net"
)

type IConnection interface {
	Start()
	Stop()
	GetTCPConnection() *net.TCPConn
	GetConnID() uint32
	RemoteAddr() net.Addr
	SendMsg(msgID uint32, data []byte) error
}

type handleFunc func(*net.TCPConn, []byte, int) error

type connection struct {
	conn         *net.TCPConn
	connID       uint32
	isClosed     bool
	router       IRouter
	exitBuffChan chan bool
}

func NewConnection(conn *net.TCPConn, connID uint32, router IRouter) *connection {
	c := &connection{
		conn:         conn,
		connID:       connID,
		isClosed:     false,
		router:       router,
		exitBuffChan: make(chan bool, 1),
	}
	return c
}

func (c *connection) startReader() {
	fmt.Println("Reader goroutine is running")
	defer fmt.Println(c.RemoteAddr().String(), "conn reader exit")
	defer c.Stop()

	for {
		pkt := NewPacket()
		headData := make([]byte, pkt.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read message head error ", err)
			c.exitBuffChan <- true
			continue
		}
		msg, err := pkt.UnPack(headData)
		if err != nil {
			fmt.Println("unpack error ", err)
			c.exitBuffChan <- true
			continue
		}
		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read message data error ", err)
				c.exitBuffChan <- true
				continue
			}
		}
		msg.SetData(data)

		req := &request{
			conn: c,
			msg: msg,
		}

		go func(request IRequest) {
			c.router.PreHandle(request)
			c.router.Handle(request)
			c.router.PostHandle(request)
		}(req)
	}
}

func (c *connection) Start() {
	go c.startReader()
	for {
		select {
		case <-c.exitBuffChan:
			return
		}
	}
}

func (c *connection) Stop() {
	if c.isClosed == true {
		return
	}
	c.isClosed = true
	c.conn.Close()
	c.exitBuffChan <- true
	close(c.exitBuffChan)
}

func (c *connection) GetTCPConnection() *net.TCPConn {
	return c.conn
}

func (c *connection) GetConnID() uint32 {
	return c.connID
}

func (c *connection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *connection) SendMsg(msgID uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection closed before send message\n")
	}
	pkt := NewPacket()
	msg, err := pkt.Pack(NewMsgPacket(msgID, data))
	if err != nil {
		fmt.Println("Pack error message ID = ", msgID)
		return errors.New("Pack error msg\n")
	}
	if _, err := c.conn.Write(msg); err != nil {
		fmt.Println("Write  message ID = ", msgID, " error", err)
		c.exitBuffChan <-true
		return errors.New("connection write error")
	}
	return nil
}
