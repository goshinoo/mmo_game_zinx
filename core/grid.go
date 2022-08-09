package core

import (
	"fmt"
	"sync"
)

// Grid 一个AOI地图中的格子类型
type Grid struct {
	//格子ID
	GID int
	//格子左边边界坐标
	MinX int
	//格子右边边界坐标
	MaxX int
	//格子上边边界坐标
	MinY int
	//格子下边边界坐标
	MaxY int
	//当前格子内玩家或物体成员的ID集合
	playerIDs map[int]struct{}
	//保护当前集合的锁
	pIDLock sync.RWMutex
}

func NewGrid(GID, MinX, MaxX, MinY, MaxY int) *Grid {
	return &Grid{
		GID:       GID,
		MinX:      MinX,
		MaxX:      MaxX,
		MinY:      MinY,
		MaxY:      MaxY,
		playerIDs: make(map[int]struct{}),
	}
}

func (g *Grid) Add(playerID int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	g.playerIDs[playerID] = struct{}{}
}

func (g *Grid) Remove(playerID int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	delete(g.playerIDs, playerID)
}

func (g *Grid) GetPlayerIDs() (playerIDs []int) {
	g.pIDLock.RLock()
	defer g.pIDLock.RUnlock()

	for k := range g.playerIDs {
		playerIDs = append(playerIDs, k)
	}
	return
}

func (g *Grid) String() string {
	return fmt.Sprintf("Grid id: %d, minX:%d, maxX:%d, minY:%d, maxY:%d, playerIDs: %v", g.GID, g.MinX, g.MaxX, g.MinY, g.MaxY, g.playerIDs)
}
