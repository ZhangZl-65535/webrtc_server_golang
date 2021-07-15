package main

import (
	"github.com/gin-gonic/gin"
	"./api"
	"./util"
	rout "./router"
)

func main() {
	configFile := "config.json"
	cf := &util.Config{}

	if err := cf.Load(configFile); err != nil {
		panic("error to load the config file:" + err.Error())
	}

	if cf.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	rout.InitRouters(router)

	api.InitCore()

	router.Run(":" + cf.Server.Port)
}
