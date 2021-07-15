package manager

import (
)

type CoreManager struct {
	DataMgr *DataManager
}

func NewCoreManager() *CoreManager {
	return &CoreManager{
		DataMgr: NewDataManager(),
	}
}

var (
	CoreMgr *CoreManager
)

func InitManagers() {
	CoreMgr = NewCoreManager()
}
