package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"../model"
)

type Room struct {
	RoomID  string `json:"roomid"`  // 房间ID
	UserID  string `json:"userid"`  // 创建者ID
	MaxSize int    `json:"maxsize"` // 最大人数
	CurSize int    `json:"cursize"` // 当前人数
}

func newRoomFromModel(mdlItem *model.Room) *Room {
	return &Room{
		RoomID:  mdlItem.GetRoomID(),
		UserID:  mdlItem.GetUserID(),
		MaxSize: mdlItem.GetMaxSize(),
		CurSize: mdlItem.GetCurSize(),
	}
}

//==== API ====
func RoomList(c *gin.Context) {
	rooms := coreMgr.DataMgr.RoomGetList()
	res := []*Room{}
	for _, room := range rooms {
		res = append(res, newRoomFromModel(room))
	}

	c.JSON(http.StatusOK, res)
}
