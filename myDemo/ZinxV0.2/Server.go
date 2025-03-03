package main

import "zinx/zinx/znet"

/*
基于Zinx框架来开发的 服务器端应用程序
*/
func main() {
	//1 创建一个server句柄,使用Zinx的api
	s := znet.NewServer("[Zinx V0.2]")
	//2 启动服务
	s.Serve()
}
