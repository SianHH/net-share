package router

import (
	"github.com/gin-gonic/gin"
	"net-share/server/api/client_tunnel"
	"net-share/server/router/middleware"
)

func setClientTunnelRouter(engine *gin.Engine) {
	api := engine.Group("api")
	v1 := api.Group("v1")
	// 端口转发
	clientTunnelGroup := v1.Group("client/tunnel", middleware.LoginAuth)
	clientTunnelGroup.GET("list", client_tunnel.Api.List)
	clientTunnelGroup.GET("query", client_tunnel.Api.Query)
	clientTunnelGroup.POST("update", client_tunnel.Api.Update)
	clientTunnelGroup.POST("create", client_tunnel.Api.Create)
	clientTunnelGroup.GET("delete", client_tunnel.Api.Delete)
}
