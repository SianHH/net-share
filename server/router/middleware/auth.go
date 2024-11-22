package middleware

import (
	"github.com/gin-gonic/gin"
	"net-share/pkg/bean"
	"net-share/pkg/jwt"
	"net-share/server/global"
)

func LoginAuth(c *gin.Context) {
	// 获取token
	token := GetToken(c)
	if token == "" {
		bean.Response.AuthNoLogin(c)
		c.Abort()
		return
	}

	claims, err := global.JwtTool.ValidToken(token)
	if err != nil {
		bean.Response.AuthInvalid(c)
		c.Abort()
		return
	}
	c.Set("claims", claims)
}

func GetToken(c *gin.Context) string {
	return c.Request.Header.Get("Token")
}

func GetClaims(c *gin.Context) jwt.Claims {
	value, ok := c.Get("claims")
	if ok {
		return value.(jwt.Claims)
	}
	return jwt.Claims{}
}
