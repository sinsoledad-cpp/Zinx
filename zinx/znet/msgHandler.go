package znet

import (
	"fmt"
	"strconv"
	"zinx/zinx/ziface"
)

/*
消息处理模块的实现
*/
type MsgHandler struct {
	//消息处理模块的API
	Apis map[uint32]ziface.IRouter
}

// 初始化/创建MsgHandle方法
func NewMsgHandle() *MsgHandler {
	return &MsgHandler{
		Apis: make(map[uint32]ziface.IRouter),
	}
}

// 调赴/执行对应的Router消息处理方法
func (mh *MsgHandler) DoMsgHandler(request ziface.IRequest) {
	//1 从Request中找到msgID
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgID = ", request.GetMsgID(), " is not found! Need Register!")
		return
	}
	//2 根据MsgID调度对应router业务即可
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// 为消息添加具体的处理逻辑
func (mh *MsgHandler) AddRouter(msgID uint32, router ziface.IRouter) {
	//1 判断 当前msg绑定的API处理方法是否已经存在
	if _, ok := mh.Apis[msgID]; ok {
		//id 已经注册了
		panic("repest api, msgID = " + strconv.Itoa(int(msgID)))
	}
	//2 添加msg与API的绑定关系
	mh.Apis[msgID] = router
	fmt.Println("Add api msgID = ", msgID, " succ!")
}
