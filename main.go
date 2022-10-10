package main

import (
	"fmt"
	"github.com/goshinoo/learn-zinx/ziface"
	"github.com/goshinoo/learn-zinx/znet"
	"mmo_game_zinx/core"
)

func main() {
	s := znet.NewServer()

	//连接创建和销毁hook钩子函数
	s.SetOnConnStart(func(connection ziface.IConnection) {
		player := core.NewPlayer(connection)
		//同步pid
		player.SyncPid()
		//广播位置
		player.BroadcastStartPosition()

		//将新上线的玩家添加到worldManager中
		core.WorldMgrObj.AddPlayer(player)

		fmt.Println("====> Player Pid = ", player.Pid, " is arrived!")
	})

	//启动服务
	s.Serve()
}
