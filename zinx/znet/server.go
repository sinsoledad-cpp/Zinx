package znet

import (
	"fmt"
	"net"
	"zinx/zinx/utils"
	"zinx/zinx/ziface"
)

// iServer 的接口实现，定义一个Server的服务器模块
type Server struct {
	// 服务器的名称
	Name string
	// 服务器绑定的IP版本
	IPVersion string
	// 服务器监听的IP地址
	IP string
	// 服务器监听的端口号
	Port int
	//当前server的消息管理模块,用来绑定MsgID和对应的处理业务API关系
	MsgHandler ziface.IMsgHandler
}

// 启动服务器方法
func (s *Server) Start() {
	fmt.Printf("[Zinx] Server Name: %s,listenner at IP: %s, Port: %d is starting \n",
		utils.GlobalObject.Name,
		utils.GlobalObject.Host,
		utils.GlobalObject.TcpPort)
	fmt.Printf("[Zinx] Version: %s,MaxConn: %d, MaxPackeeetSize: %d \n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPackageSize)
	go func() {
		//0 开启消息队列及Worker工作池
		s.MsgHandler.StartWorkerPool()
		//1 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr err: ", err)
			return
		}
		//2 监听服务器的地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVersion, "err", err)
			return
		}
		fmt.Println("start Zinx server succ, ", s.Name, " succ, Listening...")

		var cid uint32
		cid = 0
		//3 阻塞的等待客户端链接,处理客户端链接业务(读写)
		for {
			//如果有客户端链接过来，阻塞会返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err: ", err)
				continue
			}
			//将处理新连接的业务方法和conn进行绑定 得到我们的链接模块
			dealConn := NewConnection(conn, cid, s.MsgHandler)
			cid++
			//启动当前的链接业务处理
			go dealConn.Start()
		}
	}()
}

// 停止服务器方法
func (s *Server) Stop() {
	//TODO 将一些服务器的资源,状态或者一些已经开辟的链接信息,进行停止或者回收
}

// 运行服务器方法
func (s *Server) Serve() {
	// 启动server的服务功能
	s.Start()
	// TODO 做一些启动服务器之后的额外业务
	//阻塞状态
	select {}
}

// 路由功能：给当前服务注册一个路由方法，供客户端的链接处理使用
func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("Add Router succ!!")
}

/*
初始化Server模块的方法
*/
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandle(),
	}
	return s
}
