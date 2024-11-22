package bean

import (
	"github.com/gin-gonic/gin"
)

const (
	OK           = 0 // 成功
	Fail         = 1 // 失败
	AuthFail     = 2 // 登录失效
	AuthNotfound = 3 // 未登录
	AuthNotApi   = 4 // 接口未授权
)

type response struct{}

var Response = response{}

// Success 成功
func (response) Success(c *gin.Context, msg string, data interface{}) {
	result := gin.H{
		"code": OK,
	}
	if data != nil {
		result["data"] = data
	}
	if msg != "" {
		result["msg"] = msg
	}
	c.JSON(200, result)
}

// ParamFail 参数错误
func (response) ParamFail(c *gin.Context) {
	c.JSON(200, gin.H{
		"code": Fail,
		"msg":  "参数错误",
	})
}

// Fail 失败
func (response) Fail(c *gin.Context, msg string) {
	c.JSON(200, gin.H{
		"code": Fail,
		"msg":  msg,
	})
}

// AuthInvalid 登录失效
func (response) AuthInvalid(c *gin.Context) {
	c.JSON(200, gin.H{
		"code": AuthFail,
		"msg":  "登录失效",
	})
}

// AuthNoLogin 未登录
func (response) AuthNoLogin(c *gin.Context) {
	c.JSON(200, gin.H{
		"code": AuthNotfound,
		"msg":  "未登录",
	})
}

// AuthNotAdmin 不是管理员
func (response) AuthNotAdmin(c *gin.Context) {
	c.JSON(200, gin.H{
		"code": AuthNotApi,
		"msg":  "非管理员",
	})
}
