package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net-share/pkg/bean"
	"net-share/server/global"
)

// Exception 处理异常
func Exception() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				global.Logger.Error("未知错误", zap.Any("error", err))
				bean.Response.Fail(c, "系统错误")
			}
			c.Abort()
		}()
		c.Next()
	}
}
