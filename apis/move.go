package apis

import (
	"fmt"
	"github.com/goshinoo/learn-zinx/ziface"
	"github.com/goshinoo/learn-zinx/znet"
	"google.golang.org/protobuf/proto"
	"mmo_game_zinx/core"
	"mmo_game_zinx/pb"
)

// MoveApi 世界聊天 路由业务
type MoveApi struct {
	znet.BaseRouter
}

func (m *MoveApi) Handle(request ziface.IRequest) {
	//解析客户端传进来的proto协议
	proto_msg := &pb.Position{}
	err := proto.Unmarshal(request.GetData(), proto_msg)
	if err != nil {
		fmt.Println("Move unmarshal error ", err)
		return
	}

	//当前数据是属于哪个玩家发送的
	pid, err := request.GetConnection().GetProperty("pid")
	if err != nil {
		fmt.Println("Move get pid error ", err)
		return
	}

	fmt.Printf("Player pid = %d, move(%f,%f,%f,%f)\n", pid, proto_msg.X, proto_msg.Y, proto_msg.Z, proto_msg.V)

	//根据pid得到对应的player对象
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))

	player.UpdatePos(proto_msg.X, proto_msg.Y, proto_msg.Z, proto_msg.V)
}
