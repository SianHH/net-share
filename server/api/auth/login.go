package auth

import (
	"github.com/gin-gonic/gin"
	"net-share/pkg/bean"
	"net-share/server/service/auth"
)

func (*api) Login(c *gin.Context) {
	var params auth.LoginRequest
	if err := c.ShouldBindJSON(&params); err != nil {
		bean.Response.Fail(c, err.Error())
		return
	}
	token, err := auth.Service.Login(c, params)
	if err != nil {
		bean.Response.Fail(c, err.Error())
		return
	}
	bean.Response.Success(c, "", token)
}
