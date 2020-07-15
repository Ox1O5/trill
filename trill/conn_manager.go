package trill

import (
	"fmt"
	"sync"
)

type IConnManager interface {
	Add (conn IConnection)
	Remove (conn IConnection)
	Get (connID uint32) (IConnection, error)
	Len() int
	ClearConn()
}

type connManager struct {
	connections map[uint32] IConnection
	connLock sync.RWMutex
}

func NewConnManager() *connManager {
	return &connManager{
		connections: make(map[uint32] IConnection),
	}
}

func (c *connManager) Add(conn IConnection) {
	c.connLock.Lock()
	defer c.connLock.Unlock()
	c.connections[conn.GetConnID()] = conn
	fmt.Println("connection add to conn manager success : conn num = ", c.Len())
}

func (c *connManager) Remove(conn IConnection) {
	c.connLock.Lock()
	defer c.connLock.Unlock()
	delete(c.connections, conn.GetConnID())
	fmt.Println("connection Remove ConnID=",conn.GetConnID(), " successfully: conn num = ", c.Len())
}

func (c *connManager) Get(connID uint32) (IConnection, error) {
	c.connLock.RLock()
	defer c.connLock.RUnlock()
	if conn, ok := c.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, error.New("connection not found")
	}
}

func (c *connManager) Len() int {
	return len(c.connections)
}

func (c *connManager) ClearConn() {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	for connID, conn := range c.connections {
		conn.Stop()
		delete(c.connections, connID)
	}
	fmt.Println("Clear All Connections successfully: conn num = ", c.Len())
}

