package model

import (
	"errors"
	"github.com/gorilla/websocket"
	"sync"
)

const (
	DEVICE_TYPE_PHONE = 0
	DEVICE_TYPE_PC    = 1

	DEVICE_STATUS_IDLE   = 0
	DEVICE_STATUS_INCALL = 1
)
var (	
	Err_User_Not_Found = errors.New("user not found")
	Err_Conn_Not_Found = errors.New("conn not found")
)

// Helper to make Gorilla Websockets threadsafe
type threadSafeWriter struct {
	conn *websocket.Conn
	sync.Mutex
}

type User struct {
	userID string
	avatar string
	device int
	status int
	writer *threadSafeWriter
}

func NewUser(userID, avatar string, deviceType int) *User {
	return &User{
		userID: userID,
		avatar: avatar,
		device: deviceType,
		writer: &threadSafeWriter{},
	}
}

func (this *User) GetUserID() string {
	return this.userID
}
// func (this *User) SetUserID(userID string) {
// 	this.userID = userID
// }

func (this *User) GetAvatar() string {
	return this.avatar
}
func (this *User) SetAvatar(avatar string) {
	this.avatar = avatar
}

func (this *User) SetConn(conn *websocket.Conn) {
	this.writer.conn = conn
}
// func (this *User) SetConn() *websocket.Conn {
// 	return this.writer.conn
// }
func (this *User) Close() {
	this.writer.conn.Close()
	this.writer.conn = nil
}
func (this *User) WriteJSON(v interface{}) error {
	if this.writer.conn == nil {
		return Err_Conn_Not_Found
	}
	this.writer.Lock()
	defer this.writer.Unlock()

	return this.writer.conn.WriteJSON(v)
}

func (this *User) GetDevice() int {
	return this.device
}
func (this *User) SetDevice(device int) {
	this.device = device
}

func (this *User) GetStatus() int {
	return this.status
}
func (this *User) SetStatus(status int) {
	this.status = status
}