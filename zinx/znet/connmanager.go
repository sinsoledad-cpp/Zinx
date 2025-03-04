package znet

import (
	"errors"
	"fmt"
	"sync"
	"zinx/zinx/ziface"
)

/*
连接管理模块
*/
type ConnManager struct {
	connections map[uint32]ziface.IConnection //管理的连接集合
	connLock    sync.RWMutex                  //保护连接集合的读写锁
}

// 创建当前连接的方法
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

// 添加链接
func (connMgr *ConnManager) Add(conn ziface.IConnection) {
	//保护共享资源map,加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	// 将conn添加到ConnManager中
	connMgr.connections[conn.GetConnID()] = conn

	fmt.Println("connID = ", conn.GetConnID(), " add to ConnManager successfully: conn num = ", connMgr.Len())
	return
}

// 删除链接
func (connMgr *ConnManager) Remove(conn ziface.IConnection) {
	//保护共享资源map,加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//删除链接信息
	delete(connMgr.connections, conn.GetConnID())

	fmt.Println("Clear ConnID=", conn.GetConnID(), "remove from ConnManager successfully: conn num = ", connMgr.Len())
	return
}

// 根据ConnID获取链接
func (connMgr *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	//保护共享资源map,加读锁
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not found")
	}
}

// 得到当前的链接总数
func (connMgr *ConnManager) Len() int {
	return len(connMgr.connections)
}

// 清除并终止所有链接
func (connMgr *ConnManager) ClearConn() {
	//保护共享资源map,加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//删除conn并停止conn的工作
	for connID, conn := range connMgr.connections {
		//停止
		conn.Stop()
		//删除
		delete(connMgr.connections, connID)
	}
	fmt.Println("Clear All Connections successfully: conn num = ", connMgr.Len())
}
