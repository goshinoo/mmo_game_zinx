package main

import (
	"fmt"
	"github.com/goshinoo/learn-zinx/ziface"
	"github.com/goshinoo/learn-zinx/znet"
	"mmo_game_zinx/apis"
	"mmo_game_zinx/core"
)

func main() {
	s := znet.NewServer()

	//连接创建和销毁hook钩子函数
	s.SetOnConnStart(OnConnectionAdd)
	s.SetOnConnStop(OnConnectionRemove)

	//注册路由
	s.AddRouter(2, &apis.WorldChatApi{})
	s.AddRouter(3, &apis.MoveApi{})

	//启动服务
	s.Serve()
}

func OnConnectionAdd(connection ziface.IConnection) {
	player := core.NewPlayer(connection)

	//同步pid
	player.SyncPid()

	//广播位置
	player.BroadcastStartPosition()

	//将新上线的玩家添加到worldManager中
	core.WorldMgrObj.AddPlayer(player)

	//将连接绑定一个pid属性
	connection.SetProperty("pid", player.Pid)

	//同步周边玩家,告知他们当前玩家已经上线,广播当前玩家的位置信息
	player.SyncSurrounding()

	fmt.Println("====> Player Pid = ", player.Pid, " is arrived!")
}

func OnConnectionRemove(connection ziface.IConnection) {
	pid, err := connection.GetProperty("pid")
	if err != nil {
		fmt.Println("connection get pid error: ", err)
		return
	}

	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))

	//触发玩家下线的业务
	player.Offline()

	fmt.Println("====> Player Pid = ", player.Pid, " offline...!")
}
