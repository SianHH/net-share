package client_forward

import (
	"github.com/gin-gonic/gin"
	"net-share/pkg/bean"
	"net-share/server/service/client_forward"
)

func (*api) Create(c *gin.Context) {
	var param client_forward.CreateRequest
	if err := c.ShouldBindJSON(&param); err != nil {
		bean.Response.ParamFail(c)
		return
	}
	_, err := client_forward.Service.Create(param)
	if err != nil {
		bean.Response.Fail(c, err.Error())
		return
	}
	bean.Response.Success(c, "", nil)
}
