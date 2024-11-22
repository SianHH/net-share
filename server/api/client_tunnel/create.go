package client_tunnel

import (
	"github.com/gin-gonic/gin"
	"net-share/pkg/bean"
	"net-share/server/service/client_tunnel"
)

func (*api) Create(c *gin.Context) {
	var param client_tunnel.CreateRequest
	if err := c.ShouldBindJSON(&param); err != nil {
		bean.Response.ParamFail(c)
		return
	}
	_, err := client_tunnel.Service.Create(param)
	if err != nil {
		bean.Response.Fail(c, err.Error())
		return
	}
	bean.Response.Success(c, "", nil)
}
