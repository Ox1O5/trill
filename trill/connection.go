package trill

import (
	"fmt"
	"net"
)

type IConnection interface {
	Start()
	Stop()
	GetTCPConnection() *net.TCPConn
	GetConnID() uint32
	RemoteAddr() net.Addr
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
		buf := make([]byte, 512)
		_, err := c.conn.Read(buf)
		if err != nil {
			fmt.Println("read error ", err)
			c.exitBuffChan <- true
			continue
		}

		req := &request{
			conn: c,
			data: nil,
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
