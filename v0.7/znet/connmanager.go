package znet

import (
	"errors"
	"fmt"
	"github.com/suuhui/v0.7/ziface"
	"strconv"
	"sync"
)

type ConnManager struct {
	connections map[uint32]ziface.IConnection //管理的连接信息
	connLock sync.RWMutex
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

func (connMgr *ConnManager) Add(conn ziface.IConnection) {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	connMgr.connections[conn.GetConnID()] = conn
	fmt.Println("connection add to connManager successfully, connID=", conn.GetConnID())
}

func (connMgr *ConnManager) Remove(conn ziface.IConnection) {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	delete(connMgr.connections, conn.GetConnID())
	fmt.Println("delete connection from connManager successfully, connID=", conn.GetConnID())
}

func (connMgr *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	conn, ok := connMgr.connections[connID]
	if !ok {
		return nil, errors.New("connID=" + strconv.Itoa(int(connID)) + " not found.")
	}
	return conn, nil
}

func (connMgr *ConnManager) Len() int {
	return len(connMgr.connections)
}

func (connMgr *ConnManager) ClearConn() {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	for connID, conn := range connMgr.connections {
		conn.Stop()
		delete(connMgr.connections, connID)
	}
	fmt.Println("Clear all connections successfully.")
}
