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
