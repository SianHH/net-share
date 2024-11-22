package router

import (
	"github.com/gin-gonic/gin"
	"net-share/server/api/client_host"
	"net-share/server/router/middleware"
)

func setClientHostRouter(engine *gin.Engine) {
	api := engine.Group("api")
	v1 := api.Group("v1")
	// 域名解析
	clientHostGroup := v1.Group("client/host", middleware.LoginAuth)
	clientHostGroup.GET("list", client_host.Api.List)
	clientHostGroup.GET("query", client_host.Api.Query)
	clientHostGroup.POST("update", client_host.Api.Update)
	clientHostGroup.POST("create", client_host.Api.Create)
	clientHostGroup.GET("delete", client_host.Api.Delete)
}
