package router

import (
	"github.com/gin-gonic/gin"
	"net-share/server/api/auth"
	"net-share/server/router/middleware"
)

func setAuthRouter(engine *gin.Engine) {
	api := engine.Group("api")
	v1 := api.Group("v1")

	// 授权
	authGroup := v1.Group("auth")
	authGroup.POST("login", auth.Api.Login)
	authGroup.GET("captcha", auth.Api.Captcha)

	authGroup.GET("renew", middleware.LoginAuth, auth.Api.Renew)
}
