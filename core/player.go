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
		X:    float32(250 + rand.Intn(5)), //随机坐标
		Y:    0,
		Z:    float32(250 + rand.Intn(5)), //随机坐标
		V:    0,
	}
}

// SendMsg 提供一个发送给客户端的方法
//主要是将pb的protobuf序列化之后再调用zinx的Sendmsg方法
func (p *Player) SendMsg(msgId uint32, data proto.Message) {
	if p == nil {
		return
	}

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

func (p *Player) UpdatePos(X, Y, Z, V float32) {
	//旧坐标
	oldGid := WorldMgrObj.AoiManager.GetGidByPos(p.X, p.Z)

	//更新坐标
	p.X = X
	p.Z = Z
	p.V = V
	p.Y = Y
	newGid := WorldMgrObj.AoiManager.GetGidByPos(p.X, p.Z)

	//判断玩家是否跨越格子
	if newGid != oldGid {
		//新九宫格
		newgrids := WorldMgrObj.AoiManager.GetSurroundGridsByGid(newGid)
		//旧九宫格
		oldgrids := WorldMgrObj.AoiManager.GetSurroundGridsByGid(oldGid)

		//属于旧九宫格但不属于新九宫格的格子内的玩家,需要处理视野消失
		disappearGrids := make([]*Grid, 0, len(oldgrids))
		for _, o := range oldgrids {
			flg := false
			for _, n := range newgrids {
				if n == o {
					flg = true
				}
			}
			if !flg {
				disappearGrids = append(disappearGrids, o)
			}
		}
		syncPid_proto_msg := &pb.SyncPid{Pid: p.Pid}
		for _, gid := range disappearGrids {
			pids := WorldMgrObj.AoiManager.GetPidsByGid(gid.GID)
			for _, pid := range pids {
				player := WorldMgrObj.GetPlayerByPid(int32(pid))
				//发送消失消息给其他玩家
				player.SendMsg(201, syncPid_proto_msg)
				//给自己发送其他玩家消失
				p.SendMsg(201, &pb.SyncPid{Pid: int32(pid)})
			}
		}

		//属于新九宫格但不属于旧九宫格内的玩家,需要处理视野的出现
		appearGrids := make([]*Grid, 0, len(newgrids))
		for _, o := range newgrids {
			flg := false
			for _, n := range oldgrids {
				if n == o {
					flg = true
				}
			}
			if !flg {
				appearGrids = append(appearGrids, o)
			}
		}
		broadcast_proto_msg := &pb.BroadCast{
			Pid: p.Pid,
			Tp:  2,
			Data: &pb.BroadCast_P{P: &pb.Position{
				X: X,
				Y: Y,
				Z: Z,
				V: V,
			}},
		}
		for _, gid := range appearGrids {
			pids := WorldMgrObj.AoiManager.GetPidsByGid(gid.GID)
			for _, pid := range pids {
				player := WorldMgrObj.GetPlayerByPid(int32(pid))
				//发送同步消息给其他玩家
				player.SendMsg(200, broadcast_proto_msg)
				//给自己发送其他玩家出现
				p.SendMsg(200, &pb.BroadCast{
					Pid: player.Pid,
					Tp:  2,
					Data: &pb.BroadCast_P{P: &pb.Position{
						X: player.X,
						Y: player.Y,
						Z: player.Z,
						V: player.V,
					}},
				})
			}
		}

		WorldMgrObj.AoiManager.RemovePidFromGrid(int(p.Pid), oldGid)
		WorldMgrObj.AoiManager.AddPidToGrid(int(p.Pid), newGid)
	}

	//组建MsgId:200 proto数据
	proto_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  4,
		Data: &pb.BroadCast_P{P: &pb.Position{
			X: X,
			Y: Y,
			Z: Z,
			V: V,
		}},
	}

	players := p.GetSurroundingPlayers()

	for _, player := range players {
		player.SendMsg(200, proto_msg)
	}
}

// GetSurroundingPlayers 获取当前玩家周边九宫格之内的玩家
func (p *Player) GetSurroundingPlayers() []*Player {
	pids := WorldMgrObj.AoiManager.GetPidsByPos(p.X, p.Z)
	list := make([]*Player, 0, len(pids))
	for _, id := range pids {
		list = append(list, WorldMgrObj.Players[int32(id)])
	}
	return list
}

func (p *Player) Offline() {
	players := p.GetSurroundingPlayers()

	proto_msg := &pb.SyncPid{Pid: p.Pid}

	for _, player := range players {
		player.SendMsg(201, proto_msg)
	}

	WorldMgrObj.RemovePlayer(p.Pid)
}
