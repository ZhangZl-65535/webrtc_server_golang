package manager

import (
	"../model"
	"../util"
)

func (this *DataManager) InitRoom() {
	this.rooms = util.NewSafeMap()
}

func (this *DataManager) RoomGetList() []*model.Room {
	list := []*model.Room{}
	for _, item := range this.rooms.Items() {
		list = append(list, item.(*model.Room))
	}
	return list
}

func (this *DataManager) RoomGet(id string) *model.Room {
	item := this.rooms.Get(id)
	if item == nil {
		return nil
	}
	return item.(*model.Room)
}

func (this *DataManager) RoomSet(id string, room *model.Room) {
	this.rooms.Set(id, room)
}

func (this *DataManager) RoomDelete(id string) {
	this.rooms.Delete(id)
}
