package middleware

import "github.com/gin-gonic/gin"

// Cache @Title 设置静态文件缓存
func Cache() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Cache-Control", "public,max-age=86400")
	}
}

// NoCache @Title 不缓存
func NoCache() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Cache-Control", "no-cache")
	}
}
