package manager

import (
	"../util"
)

type DataManager struct {
	rooms *util.SafeMap
	users *util.SafeMap
}

func NewDataManager() *DataManager {
	m := &DataManager{
	}
	m.InitRoom()
	m.InitUser()
	return m
}
