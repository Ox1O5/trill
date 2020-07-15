package trill

import (
	"errors"
	"fmt"
	"github.com/Ox1O5/trill/utils"
	"io"
	"net"
	"sync"
)

type IConnection interface {
	Start()
	Stop()
	GetTCPConnection() *net.TCPConn
	GetConnID() uint32
	RemoteAddr() net.Addr
	SendMsg(msgID uint32, data []byte) error
	SendBuffMsg(msgID uint32, data []byte) error
	SetProperty(key string, value interface{})
	GetProperty(key string)(interface{}, error)
	RemoveProperty(key string)
}

type handleFunc func(*net.TCPConn, []byte, int) error

type connection struct {
	TcpServer    IServer
	conn         *net.TCPConn
	connID       uint32
	isClosed     bool
	msgHandler   IMsgHandle
	exitBuffChan chan bool
	msgChan      chan []byte
	msgBuffChan  chan []byte
	property     map[string]interface{}
	propertyLock sync.RWMutex
}

func NewConnection(server IServer, conn *net.TCPConn, connID uint32, msgHandler IMsgHandle) *connection {
	c := &connection{
		TcpServer:    server,
		conn:         conn,
		connID:       connID,
		isClosed:     false,
		msgHandler:   msgHandler,
		exitBuffChan: make(chan bool, 1),
		msgChan:      make(chan []byte),
		msgBuffChan:  make(chan []byte, utils.GlobalObject.MaxMsgChanLen),
		property:     make(map[string]interface{}),
	}
	c.TcpServer.GetConnManager().Add(c)
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
			fmt.Println("read Message head error ", err)
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
				fmt.Println("read Message Data error ", err)
				c.exitBuffChan <- true
				continue
			}
		}
		msg.SetData(data)

		req := &request{
			conn: c,
			msg:  msg,
		}
		if utils.GlobalObject.WorkerPoolSize > 0 {
			c.msgHandler.SendMsgToTaskQueue(req)
		} else {
			go c.msgHandler.DoMsgHandle(req)
		}
	}
}

func (c *connection) startWriter() {
	fmt.Println("[Writer goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn writer exit]")

	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.conn.Write(data); err != nil {
				fmt.Println("Send Data error: ", err, " conn writer exit")
				return
			}
		case data, ok := <-c.msgBuffChan:
			if ok {
				if _, err := c.conn.Write(data); err != nil {
					fmt.Println("Send Data error: ", err, " conn writer exit")
					return
				}
			} else {
				fmt.Println("msgBuffChan is Closed")
				break
			}

		case <-c.exitBuffChan:
			return
		}
	}
}

func (c *connection) Start() {
	go c.startReader()
	go c.startWriter()
	c.TcpServer.CallOnConnStart(c)
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
	c.TcpServer.CallOnConnStop(c)
	c.conn.Close()
	c.exitBuffChan <- true
	c.TcpServer.GetConnManager().Remove(c)
	close(c.exitBuffChan)
	close(c.msgBuffChan)
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
		return errors.New("Connection closed before send Message\n")
	}
	pkt := NewPacket()
	msg, err := pkt.Pack(NewMsgPacket(msgID, data))
	if err != nil {
		fmt.Println("Pack error Message ID = ", msgID)
		return errors.New("Pack error msg\n")
	}
	c.msgChan <- msg
	return nil
}

func (c *connection) SendBuffMsg(msgID uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection closed before send Message\n")
	}
	pkt := NewPacket()
	msg, err := pkt.Pack(NewMsgPacket(msgID, data))
	if err != nil {
		fmt.Println("Pack error Message ID = ", msgID)
		return errors.New("Pack error msg\n")
	}
	c.msgBuffChan <- msg
	return nil
}

func (c *connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	c.property[key] = value
}
//获取链接属性
func (c *connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	if value, ok := c.property[key]; ok  {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}
//移除链接属性
func (c *connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	delete(c.property, key)
}