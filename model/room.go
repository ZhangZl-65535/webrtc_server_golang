package model

import (
	"../util"
)

type Room struct {
	roomID  string        // 房间ID
	userID  string        // 创建者ID
	maxSize int           // 最大人数
	users   *util.SafeMap // 房间内成员
}

func NewRoom(roomID, userID string, maxSize int) *Room {
	return &Room{
		roomID:  roomID,
		userID:  userID,
		maxSize: maxSize,
		users:   util.NewSafeMap(),
	}
}

func (this *Room) GetRoomID() string {
	return this.roomID
}
// func (this *Room) SetRoomID(roomID string) {
// 	this.roomID = roomID
// }

func (this *Room) GetUserID() string {
	return this.userID
}
// func (this *Room) SetUserID(userID string) {
// 	this.userID = userID
// }

func (this *Room) GetMaxSize() int {
	return this.maxSize
}
func (this *Room) SetMaxSize(maxSize int) {
	this.maxSize = maxSize
}

func (this *Room) AddUser(user *User) {
	this.users.Set(user.GetUserID(), user)
}
func (this *Room) RemoveUser(userID string) {
	this.users.Delete(userID)
}
func (this *Room) GetCurSize() int {
	return this.users.Length()
}

func (this *Room) GetUsers() []*User {
	users := []*User{}
	for _, item := range this.users.Items() {
		if item == nil {
			continue
		}
		users = append(users, item.(*User))
	}
	return users
}
func (this *Room) GetUserIDs() []string {
	users := []string{}
	for _, item := range this.users.Items() {
		if item == nil {
			continue
		}
		users = append(users, item.(*User).GetUserID())
	}
	return users
}