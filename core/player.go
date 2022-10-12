package core

import (
	"fmt"
	"github.com/goshinoo/learn-zinx/ziface"
	"google.golang.org/protobuf/proto"
	"math/rand"
	"mmo_game_zinx/pb"
	"sync"
)

var PidGen int32 = 1
var IdLock sync.Mutex

type Player struct {
	Pid  int32              // 玩家id
	Conn ziface.IConnection //客户端连接
	X    float32            //平面x坐标
	Y    float32            //高度
	Z    float32            //平面y坐标
	V    float32            //旋转的角度0-360
}

func NewPlayer(conn ziface.IConnection) *Player {
	//生成玩家id
	IdLock.Lock()
	id := PidGen
	PidGen++
	IdLock.Unlock()

	return &Player{
		Pid:  id,
		Conn: conn,
		X:    float32(160 + rand.Intn(10)), //随机坐标
		Y:    0,
		Z:    float32(140 + rand.Intn(20)), //随机坐标
		V:    0,
	}
}

// SendMsg 提供一个发送给客户端的方法
//主要是将pb的protobuf序列化之后再调用zinx的Sendmsg方法
func (p *Player) SendMsg(msgId uint32, data proto.Message) {
	msg, err := proto.Marshal(data)
	if err != nil {
		fmt.Println("proto marshal err:", err)
		return
	}

	if p.Conn == nil {
		fmt.Println("conn in player is closed")
		return
	}
	err = p.Conn.SendMsg(msgId, msg)
	if err != nil {
		fmt.Println("send msg err:", err)
	}
}

func (p *Player) SyncPid() {
	proto_msg := &pb.SyncPid{Pid: p.Pid}
	p.SendMsg(1, proto_msg)
}

func (p *Player) BroadcastStartPosition() {
	proto_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}
	p.SendMsg(200, proto_msg)

}

// Talk 玩家广播世界聊天消息
func (p *Player) Talk(content string) {
	//组建MsgId:200 proto数据
	proto_msg := &pb.BroadCast{
		Pid:  p.Pid,
		Tp:   1,
		Data: &pb.BroadCast_Content{Content: content},
	}
	//得到当前世界所有在线玩家
	players := WorldMgrObj.GetAllPlayers()

	//向所有玩家(包括自己)发送200消息
	for _, player := range players {
		player.SendMsg(200, proto_msg)
	}
}

func (p *Player) SyncSurrounding() {
	//获取当前玩家周围有哪些玩家
	pids := WorldMgrObj.AoiManager.GetPidsByPos(p.X, p.Z)
	players := make([]*Player, 0, len(pids))
	for _, pid := range pids {
		players = append(players, WorldMgrObj.GetPlayerByPid(int32(pid)))
	}
	//将当前玩家位置信息通过MsgID:200发给周围玩家(让其他玩家看到自己)
	proto_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2,
		Data: &pb.BroadCast_P{P: &pb.Position{
			X: p.X,
			Y: p.Y,
			Z: p.Z,
			V: p.V,
		}},
	}
	for _, player := range players {
		player.SendMsg(200, proto_msg)
	}

	//将周围全部玩家的位置信息发给当前玩家MsgId:202(让自己看到其他玩家)
	//制作msgID:202数据
	//制作pb.Player slice
	players_proto_msg := make([]*pb.Player, 0, len(players))
	for _, player := range players {
		players_proto_msg = append(players_proto_msg, &pb.Player{
			Pid: player.Pid,
			P: &pb.Position{
				X: player.X,
				Y: player.Y,
				Z: player.Z,
				V: player.V,
			},
		})
	}

	//封装
	syncPlayers_proto_msg := &pb.SyncPlayers{Ps: players_proto_msg}

	//将组建好的数据发送给当前玩家客户端
	p.SendMsg(202, syncPlayers_proto_msg)
}
