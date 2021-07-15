package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"../model"
)

type User struct {
	UserID  string `json:"userid"`  // 用户ID
	Avatar  string `json:"avatar"`  // 头像
	IsPhone bool   `json:"isPhone"` // 是否手机端
}

func newUserFromModel(mdlItem *model.User) *User {
	isPhone := true
	if mdlItem.GetDevice() > 0 {
		isPhone = false
	}
	return &User{
		UserID:  mdlItem.GetUserID(),
		Avatar:  mdlItem.GetAvatar(),
		IsPhone: isPhone,
	}
}

//==== API ====
func UserList(c *gin.Context) {
	users := coreMgr.DataMgr.UserGetList()
	res   := []*User{}
	for _, User := range users {
		res = append(res, newUserFromModel(User))
	}

	c.JSON(http.StatusOK, res)
}
