package manager

import (
	"../model"
	"../util"
)

func (this *DataManager) InitUser() {
	this.users = util.NewSafeMap()
}

func (this *DataManager) UserGetList() []*model.User {
	list := []*model.User{}
	for _, item := range this.users.Items() {
		list = append(list, item.(*model.User))
	}
	return list
} 

func (this *DataManager) UserGet(id string) *model.User {
	item := this.users.Get(id)
	if item == nil {
		return nil
	}
	return item.(*model.User)
}

func (this *DataManager) UserSet(id string, user *model.User) {
	this.users.Set(id, user)
}

func (this *DataManager) UserDelete(id string) {
	user := this.UserGet(id)
	if user != nil {
		user.Close()
	}
	this.users.Delete(id)
}
