package api

import (
	"../manager"
)

var (
	coreMgr *manager.CoreManager
)

func InitCore() {
	manager.InitManagers()
	coreMgr = manager.CoreMgr
}
