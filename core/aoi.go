package core

import "fmt"

const (
	AOI_MIN_X  = 0
	AOI_MAX_X  = 500
	AOI_CNTS_X = 100
	AOI_MIN_Y  = 0
	AOI_MAX_Y  = 500
	AOI_CNTS_Y = 100
)

// AOIManager AOI区域管理模块
type AOIManager struct {
	//区域左边界坐标
	MinX int
	//区域右边界坐标
	MaxX int
	//X方向格子数量
	CntsX int
	//区域上边界坐标
	MinY int
	//区域下边界坐标
	MaxY int
	//Y方向格子数量
	CntsY int
	//当前区域中有哪些格子
	grids map[int]*Grid
}

func NewAOIManager(minX, maxX, cntsX, minY, maxY, cntsY int) *AOIManager {
	aoiMgr := &AOIManager{
		MinX:  minX,
		MaxX:  maxX,
		CntsX: cntsX,
		MinY:  minY,
		MaxY:  maxY,
		CntsY: cntsY,
		grids: make(map[int]*Grid),
	}

	//给AOI初始化区域的所有格子进行编号和初始化
	for y := 0; y < cntsY; y++ {
		for x := 0; x < cntsX; x++ {
			//计算格子ID 根据x,y编号
			//格子编号:id = idy*cntsx+idx
			gid := y*cntsX + x
			aoiMgr.grids[gid] = NewGrid(
				gid,
				aoiMgr.MinX+x*aoiMgr.gridWidth(),
				aoiMgr.MinX+(x+1)*aoiMgr.gridWidth(),
				aoiMgr.MinY+y*aoiMgr.gridHeight(),
				aoiMgr.MinY+(y+1)*aoiMgr.gridHeight(),
			)
		}
	}
	return aoiMgr
}

//得到每个格子在X轴方向的宽度
func (m *AOIManager) gridWidth() int {
	return (m.MaxX - m.MinX) / m.CntsX
}

//得到每个格子在Y轴方向的宽度
func (m *AOIManager) gridHeight() int {
	return (m.MaxY - m.MinY) / m.CntsY
}

func (m *AOIManager) String() string {
	s := fmt.Sprintf("AOIManager: MinX:%d, MaxX:%d, cntsX:%d, minY:%d, maxY:%d, cntsY:%d \n", m.MinX, m.MaxX, m.CntsX, m.MinY, m.MaxY, m.CntsY)
	for _, grid := range m.grids {
		s += fmt.Sprintln(grid)
	}
	return s
}

// GetSurroundGridsByGid 通过GID获取九宫格范围
func (m *AOIManager) GetSurroundGridsByGid(GID int) (grids []*Grid) {
	//判断gid是否在AOIManager中
	if _, ok := m.grids[GID]; !ok {
		return
	}

	//将当前gid加入九宫格切片中
	grids = append(grids, m.grids[GID])

	//gid左边是否有格子?右边是否有格子?
	//需要通过gid得到当前格子x轴的编号 idx=id%nx
	idx := GID % m.CntsX

	//判断idx编号左边有格子
	if idx > 0 {
		grids = append(grids, m.grids[GID-1])
	}

	//判断idx编号右边边有格子
	if idx < m.CntsX-1 {
		grids = append(grids, m.grids[GID+1])
	}

	//将x轴当前的格子都取出.进行遍历,再分别得到每个格子上下是否还有格子
	gridsX := make([]int, 0, len(grids))
	for _, v := range grids {
		gridsX = append(gridsX, v.GID)
	}

	//遍历gridsX集合中的每个格子
	for _, v := range gridsX {
		//得到y轴编号
		idy := v / m.CntsX
		//gid上边是否还有格子
		if idy > 0 {
			grids = append(grids, m.grids[v-m.CntsX])
		}
		//gid下边是否还有格子
		if idy < m.CntsY-1 {
			grids = append(grids, m.grids[v+m.CntsX])
		}
	}

	return
}

// GetGidByPos 通过坐标得到gid
func (m *AOIManager) GetGidByPos(x, y float32) int {
	idx := (int(x) - m.MinX) / m.gridWidth()
	idy := (int(y) - m.MinY) / m.gridHeight()
	return idy*m.CntsX + idx

}

// GetPidsByPos 通过横纵坐标得到九宫格内全部playerIDs
func (m *AOIManager) GetPidsByPos(x, y float32) (playerIDs []int) {
	//得到当前玩家的GID格子
	gid := m.GetGidByPos(x, y)
	//通过GID得到九宫格信息
	grids := m.GetSurroundGridsByGid(gid)
	//将九宫格的信息全部player的id累加到playerIDs
	for _, grid := range grids {
		playerIDs = append(playerIDs, grid.GetPlayerIDs()...)
		fmt.Printf("=== grid ID: %d,pids :%v ======\n", grid.GID, grid.GetPlayerIDs())
	}

	return
}

// AddPidToGrid 添加一个playerID到格子中
func (m *AOIManager) AddPidToGrid(pID, gID int) {
	m.grids[gID].Add(pID)
}

// RemovePidFromGrid 移除格子中的playerID
func (m AOIManager) RemovePidFromGrid(pID, gID int) {
	m.grids[gID].Remove(pID)
}

// GetPidsByGid 通过GID获取全部的playerIDs
func (m *AOIManager) GetPidsByGid(gID int) []int {
	return m.grids[gID].GetPlayerIDs()
}

// AddPidToGridByPos 通过坐标添加playerID
func (m *AOIManager) AddPidToGridByPos(x, y float32, pID int) {
	m.grids[m.GetGidByPos(x, y)].Add(pID)
}

// RemovePidFromGridByPos 通过坐标删除playerID
func (m *AOIManager) RemovePidFromGridByPos(x, y float32, pID int) {
	m.grids[m.GetGidByPos(x, y)].Remove(pID)
}
