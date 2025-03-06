package main

import (
	"fmt"
	"zinx/mmo_game_zinx/core"
	"zinx/zinx/ziface"

	"zinx/zinx/znet"
)

// 当前客户端建立连接之后的hook函数
func OnConnectionAdd(conn ziface.IConnection) {
	//创建一个Player对象
	player := core.NewPlayer(conn)
	//给客户端发送MsgID:1的消息:同步当前Player的ID给客户端
	player.SyncPid()
	//给客户端发送MsgID:200的消息:同步当前Player的初始位置给客户端
	player.BroadCastStartPosition()
	fmt.Println("===>>Player pid = ", player.Pid, " is arrived <<===")
}

func main() {
	//创建服务器句柄
	s := znet.NewServer("zinx")

	//链接创建和销毁的HOOK钩子函数
	s.SetOnConnStart(OnConnectionAdd)

	//注册一些路由业务
	//// Add LTV data format Decoder
	//s.SetDecoder(zdecoder.NewLTV_Little_Decoder())
	//// Add LTV data format Pack packet Encoder
	//s.SetPacket(zpack.NewDataPackLtv())
	//启动服务
	s.Serve()
}
