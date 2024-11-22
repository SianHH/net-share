package auth

import (
	"github.com/gin-gonic/gin"
	"net-share/pkg/bean"
	"net-share/server/constant"
	"net-share/server/global"
	"net-share/server/service/auth"
)

func (*api) Captcha(c *gin.Context) {
	key, bs64, err := auth.Service.Captcha()
	if err != nil {
		bean.Response.Fail(c, err.Error())
		return
	}
	bean.Response.Success(c, "", gin.H{
		"key":      key,
		"bs64":     bs64,
		"security": global.Cache.GetString(constant.CacheSecurityIPKey+c.ClientIP()) == "",
	})
}
