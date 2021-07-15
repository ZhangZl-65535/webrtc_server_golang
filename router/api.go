package router

import (
	"github.com/gin-gonic/gin"
	"../api"
)

func initApi(router *gin.Engine) {
	apiRouter := router.Group("/")
	// room
	apiRouter.GET("roomList", api.RoomList)
	apiRouter.POST("roomList", api.RoomList)

	// user
	apiRouter.GET("userList", api.UserList)
	apiRouter.POST("userList", api.UserList)

	// ws
	wsRouter := router.Group("/ws")
	wsRouter.GET("/:userid/:device", api.WS)
}
