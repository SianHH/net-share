package client_host

import (
	"github.com/gin-gonic/gin"
	"net-share/pkg/bean"
	"net-share/server/service/client_host"
)

func (*api) Delete(c *gin.Context) {
	var param = c.Query("code")
	if param == "" {
		bean.Response.ParamFail(c)
		return
	}
	if err := client_host.Service.Delete(param); err != nil {
		bean.Response.Fail(c, err.Error())
		return
	}
	bean.Response.Success(c, "", nil)
}
