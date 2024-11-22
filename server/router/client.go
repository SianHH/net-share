package router

import (
	"github.com/gin-gonic/gin"
	"net-share/server/api/client"
	"net-share/server/router/middleware"
)

func setClientRouter(engine *gin.Engine) {
	api := engine.Group("api")
	v1 := api.Group("v1")
	// 客户端
	v1.Group("client").Any("ws", client.Api.Ws)
	clientGroup := v1.Group("client", middleware.LoginAuth)
	clientGroup.GET("list", client.Api.List)
	clientGroup.POST("update", client.Api.Update)
	clientGroup.POST("create", client.Api.Create)
	clientGroup.GET("delete", client.Api.Delete)
}
