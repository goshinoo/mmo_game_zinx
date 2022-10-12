package apis

import (
	"fmt"
	"github.com/goshinoo/learn-zinx/ziface"
	"github.com/goshinoo/learn-zinx/znet"
	"google.golang.org/protobuf/proto"
	"mmo_game_zinx/core"
	"mmo_game_zinx/pb"
)

// WorldChatApi 世界聊天 路由业务
type WorldChatApi struct {
	znet.BaseRouter
}

func (wc *WorldChatApi) Handle(request ziface.IRequest) {
	//解析客户端传进来的proto协议
	proto_msg := &pb.Talk{}
	err := proto.Unmarshal(request.GetData(), proto_msg)
	if err != nil {
		fmt.Println("Talk unmarshal error ", err)
		return
	}

	//当前聊天数据是属于哪个玩家发送的
	pid, err := request.GetConnection().GetProperty("pid")
	if err != nil {
		fmt.Println("talk get pid error ", err)
		return
	}

	//根据pid得到对应的player对象
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))

	//将消息广播给其他全部在线玩家
	player.Talk(proto_msg.Content)
}
