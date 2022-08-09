package core

import (
	"fmt"
	"testing"
)

func TestNewAOIManager(t *testing.T) {
	//初始化AOIManager
	aoiMgr := NewAOIManager(0, 250, 5, 0, 250, 6)
	fmt.Println(aoiMgr)
}

func TestAOIManager_GetSurroundGridsByGid(t *testing.T) {
	aoiMgr := NewAOIManager(0, 250, 5, 0, 300, 6)
	for gid := range aoiMgr.grids {
		//得到当前gid周边九宫格信息
		grids := aoiMgr.GetSurroundGridsByGid(gid)
		fmt.Println("gid = ", gid, ", grids len = ", len(grids))
		var gids []int
		for _, grid := range grids {
			gids = append(gids, grid.GID)
		}

		fmt.Println(gids)
	}
}
