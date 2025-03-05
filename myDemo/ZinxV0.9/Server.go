package main

import (
	"fmt"
	"zinx/zinx/ziface"
	"zinx/zinx/znet"
)

/*
基于Zinx框架来开发的 服务器端应用程序
*/
// ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

// Test Handle
func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call PingRouter Handle...")
	// 先读取客户端的数据，再回写ping...ping...ping...
	fmt.Println("recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))
	err := request.GetConnection().SendMsg(200, []byte("ping...ping...ping..."))
	if err != nil {
		fmt.Println(err)
	}
}

// hello Zinx test 自定义路由
type HelloZinxRouter struct {
	znet.BaseRouter
}

// Test Handle
func (this *HelloZinxRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call HelloZinxRouter Handle...")
	// 先读取客户端的数据，再回写ping...ping...ping...
	fmt.Println("recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))
	err := request.GetConnection().SendMsg(201, []byte("Hello Welcome to Zinx!!!"))
	if err != nil {
		fmt.Println(err)
	}
}

// 创建链接之后执行钩子函数
func DoConnectionBegin(conn ziface.IConnection) {
	fmt.Println("DoConnectionBegin is Called ... ")
	if err := conn.SendMsg(200, []byte("DoConnectionBegin...")); err != nil {
		fmt.Println(err)
	}
	//设置一些链接属性
	fmt.Println("Set conn property...")
	conn.SetProperty("Name", "姓名")
	conn.SetProperty("GithHub", "仓库")
	conn.SetProperty("Home", "家")
	conn.SetProperty("Blog", "博客")
}

// 链接断开之前的需要执行的钩子函数
func DoConnectionLost(conn ziface.IConnection) {
	fmt.Println("DoConnectionLost is Called ... ")
	fmt.Println("conn Id = ", conn.GetConnID(), " is Lost ... ")
	//获取链接属性
	if value, err := conn.GetProperty("Name"); err == nil {
		fmt.Println("Name = ", value)
	}
	if value, err := conn.GetProperty("GithHub"); err == nil {
		fmt.Println("GithHub = ", value)
	}
	if value, err := conn.GetProperty("Home"); err == nil {
		fmt.Println("Home = ", value)
	}
	if value, err := conn.GetProperty("Blog"); err == nil {
		fmt.Println("Blog = ", value)
	}
}
func main() {
	//1 创建一个server句柄,使用Zinx的api
	s := znet.NewServer("[Zinx V0.2]")
	//2 注册链接Hook钩子函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	//3 给当前zinx框架添加一个自定义的router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloZinxRouter{})
	//4 启动服务
	s.Serve()
}
