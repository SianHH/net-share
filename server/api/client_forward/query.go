package client_forward

import (
	"github.com/gin-gonic/gin"
	"net-share/pkg/bean"
	"net-share/server/service/client_forward"
)

func (*api) Query(c *gin.Context) {
	result, err := client_forward.Service.Query(c.Query("code"))
	if err != nil {
		bean.Response.Fail(c, err.Error())
		return
	}
	bean.Response.Success(c, "", result)
}
