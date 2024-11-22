package router

import (
	"github.com/gin-gonic/gin"
	"net-share/server/api/client_forward"
	"net-share/server/router/middleware"
)

func setClientForwardRouter(engine *gin.Engine) {
	api := engine.Group("api")
	v1 := api.Group("v1")
	// 端口转发
	clientForwardGroup := v1.Group("client/forward", middleware.LoginAuth)
	clientForwardGroup.GET("list", client_forward.Api.List)
	clientForwardGroup.GET("query", client_forward.Api.Query)
	clientForwardGroup.POST("update", client_forward.Api.Update)
	clientForwardGroup.POST("create", client_forward.Api.Create)
	clientForwardGroup.GET("delete", client_forward.Api.Delete)
}
