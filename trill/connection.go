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
	msgHandler IMsgHandle
	exitBuffChan chan bool
	msgChan chan []byte
}

func NewConnection(conn *net.TCPConn, connID uint32, msgHandler IMsgHandle) *connection {
	c := &connection{
		conn:         conn,
		connID:       connID,
		isClosed:     false,
		msgHandler : msgHandler,
		exitBuffChan: make(chan bool, 1),
		msgChan: make(chan []byte),
	}
	return c
}

func (c *connection) startReader() {
	fmt.Println("[Reader goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn reader exit]")
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

		go c.msgHandler.DoMsgHandle(req)
	}
}

func(c *connection)startWriter() {
	fmt.Println("[Writer goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn writer exit]")

	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.conn.Write(data); err != nil {
				fmt.Println("Send data error: ", err, " conn writer exit")
			}
		case <-c.exitBuffChan:
			return
		}
	}
}

func (c *connection) Start() {
	go c.startReader()
	go c.startWriter()
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
	c.msgChan <- msg
	return nil
}
