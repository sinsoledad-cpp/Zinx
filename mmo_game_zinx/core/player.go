package core

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"math/rand"
	"sync"
	"zinx/mmo_game_zinx/pb"
	"zinx/zinx/ziface"
)

// 玩家对象
type Player struct {
	Pid  int32              // 玩家ID
	Conn ziface.IConnection // 当前玩家的连接（用于和客户端的连接）
	X    float32            //平面的x坐标
	Y    float32            //高度
	Z    float32            //平面的y坐标（注意不是Y）
	V    float32            //旋转的0-360角度
}

/*
Player ID 生成器
*/
var PidGen int32 = 1  //用来生产玩家ID的计数器
var IdLock sync.Mutex //保护PidGen的Mutex

// 创建一个玩家的方法
func NewPlayer(conn ziface.IConnection) *Player {
	//生成要给玩家ID
	IdLock.Lock()
	id := PidGen
	PidGen++
	IdLock.Unlock()

	//创建一个玩家对象
	p := &Player{
		Pid:  id,
		Conn: conn,
		X:    float32(160 + rand.Intn(10)), //随机在160坐标点 基于X轴若干偏移
		Y:    0,
		Z:    float32(140 + rand.Intn(20)), //随机在120坐标点 基于Y轴若干偏移
		V:    0,                            //角度为0
	}
	return p
}

/*
提供一个发送个客户端消息的方法
主要是将pb的protobuf数据序列化之后，再调用zinx的SendMsg方法
*/
func (p *Player) SendMsg(msgId uint32, data proto.Message) {
	//将pb序列化
	msg, err := proto.Marshal(data)
	if err != nil {
		fmt.Println("marshal msg err:", err)
		return
	}
	//将二进制文件 通过zinx框架的sendmsg将数据发送给客户端
	if p.Conn == nil {
		fmt.Println("connection in player is nil")
		return
	}
	//调用zinx的SendMsg方法
	if err := p.Conn.SendMsg(msgId, msg); err != nil {
		fmt.Println("Player SendMsg error!")
		return
	}
	return
}

// 告知客户端玩家Pid，同步已经生成的玩家ID给客户端
func (p *Player) SyncPid() {
	//组建MsgID:0 的proto数据
	proto_msg := &pb.SyncPid{
		Pid: p.Pid,
	}
	//将消息发送给客户端
	p.SendMsg(1, proto_msg)
}

// 广播玩家自己的出生地点
func (p *Player) BroadCastStartPosition() {
	//组建MsgID:200 的proto数据
	proto_msg := &pb.Broadcast{
		Pid: p.Pid,
		Tp:  2, //2-玩家位置
		Data: &pb.Broadcast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}
	//将消息发送给客户端
	p.SendMsg(200, proto_msg)
}

// 玩家广播世界聊天消息
func (p *Player) Talk(content string) {
	//1 组建MsgID:200 proto数据
	proto_msg := &pb.Broadcast{
		Pid: p.Pid,
		Tp:  1, //tp-1 代表聊天广播
		Data: &pb.Broadcast_Content{
			Content: content,
		},
	}

	//2 得到当前世界所有的在线玩家
	players := WorldMgrObj.GetAllPlayers()

	//3 向所有的玩家(包括自己)发送MsgID:200消息
	for _, player := range players {
		//player分别给对应的客户端发送消息
		player.SendMsg(200, proto_msg)
	}

}

// 同步玩家上线的位置消息
func (p *Player) SyncSurrounding() {
	// 1 获取当前玩家周围的玩家有哪些(九宫格)
	pids := WorldMgrObj.AOIMgr.GetPidsByPos(p.X, p.Z)
	players := make([]*Player, 0, len(pids))
	for _, pid := range pids {
		players = append(players, WorldMgrObj.GetPlayerByPid(int32(pid)))
	}
	// 2 将当前玩家的位置信息通过MsgID:200 发给周围的玩家(让其他玩家看到自己)
	//2.1 组建MsgID:200 的proto数据
	proto_msg := &pb.Broadcast{
		Pid: p.Pid,
		Tp:  2, //2-广播坐标
		Data: &pb.Broadcast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}
	//2.2 全部周围的玩家都向格子的客户端发送200消息,proto_msg
	for _, player := range players {
		player.SendMsg(200, proto_msg)
	}
	// 3 将周围的全部玩家的位置信息发送给当前玩家MsgID:202客户端(让自己看到周围玩家)
	// 3.1 组建MsgID:202 proto数据
	// 3.1.1 制作pb.Player slice
	players_proto_msg := make([]*pb.Player, 0, len(players))
	for _, player := range players {
		//制作一个message Player
		p := &pb.Player{
			Pid: player.Pid,
			P: &pb.Position{
				X: player.X,
				Y: player.Y,
				Z: player.Z,
				V: player.V,
			},
		}
		players_proto_msg = append(players_proto_msg, p)
	}
	// 3.1.2 封装SyncPlayer protobuf数据
	SyncPlayers_proto_msg := &pb.SyncPlayers{
		Ps: players_proto_msg[:], //!!!注意!!!
	}
	// 3.2 将组建好的数据发送给当前玩家的客户端
	p.SendMsg(202, SyncPlayers_proto_msg)
}

// 广播当前玩家的位置移动信息
func (p *Player) UpdatePos(x, y, z, v float32) {
	// 更新当前玩家player对象的坐标
	p.X = x
	p.Y = y
	p.Z = z
	p.V = v
	//组建广播proto协议 MsgID:200 Tp-4
	proto_msg := &pb.Broadcast{
		Pid: p.Pid,
		Tp:  4, //4-广播坐标
		Data: &pb.Broadcast_P{
			P: &pb.Position{
				X: x,
				Y: y,
				Z: z,
				V: v,
			},
		},
	}
	//获取当前玩家的周边玩家AOI九宫格之内的玩家
	players := p.GetSuroundingPlayers()
	//一次给每个玩家对应的客户端发送当前玩家位置更新的消息
	for _, player := range players {
		player.SendMsg(200, proto_msg)
	}
}

// 获取当前玩家的周边玩家AOI九宫格之内的玩家
func (p *Player) GetSuroundingPlayers() []*Player {
	// 得到当前AOI九宫格内的所有玩家PID
	pids := WorldMgrObj.AOIMgr.GetPidsByPos(p.X, p.Z)
	//将所有的pid对应的Player放到Players切片中
	players := make([]*Player, 0, len(pids))
	for _, pid := range pids {
		players = append(players, WorldMgrObj.GetPlayerByPid(int32(pid)))
	}
	return players
}

// 玩家下线
func (p *Player) Offline() {
	//得到当前玩家周边的九宫格内的都有哪些玩家
	players := p.GetSuroundingPlayers()
	//给周围玩家广播MsgID:201消息
	proto_msg := &pb.SyncPid{
		Pid: p.Pid,
	}
	for _, player := range players {
		player.SendMsg(201, proto_msg)
	}
	WorldMgrObj.AOIMgr.RemoveFromGridByPos(int(p.Pid), p.X, p.Z)
	WorldMgrObj.RemovePlayerByPid(p.Pid)
}
