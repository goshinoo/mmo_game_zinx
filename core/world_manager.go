package core

import "sync"

// WorldManager 当前世界总管理模块
type WorldManager struct {
	AoiManager *AOIManager
	Players    map[int32]*Player
	pLock      sync.RWMutex
}

// WorldMgrObj 提供对外的世界观里模块的句柄
var WorldMgrObj *WorldManager

func init() {
	WorldMgrObj = &WorldManager{
		AoiManager: NewAOIManager(AOI_MIN_X, AOI_MAX_X, AOI_CNTS_X, AOI_MIN_Y, AOI_MAX_Y, AOI_CNTS_Y),
		Players:    make(map[int32]*Player),
	}
}

// AddPlayer 添加玩家
func (wm *WorldManager) AddPlayer(player *Player) {
	wm.pLock.Lock()
	wm.Players[player.Pid] = player
	wm.pLock.Unlock()

	wm.AoiManager.AddPidToGridByPos(player.X, player.Z, int(player.Pid))
}

// RemovePlayer 删除玩家
func (wm *WorldManager) RemovePlayer(pid int32) {
	wm.pLock.RLock()
	player := wm.Players[pid]
	wm.pLock.RUnlock()

	wm.AoiManager.RemovePidFromGridByPos(player.X, player.Z, int(pid))

	wm.pLock.Lock()
	delete(wm.Players, player.Pid)
	wm.pLock.Unlock()
}

// GetPlayerByPid 通过玩家ID查询player对象
func (wm *WorldManager) GetPlayerByPid(pid int32) *Player {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()
	return wm.Players[pid]
}

func (wm *WorldManager) GetAllPlayers() []*Player {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()

	players := make([]*Player, 0)

	for _, v := range wm.Players {
		players = append(players, v)
	}

	return players
}
