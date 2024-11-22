package client_host

import (
	"github.com/gin-gonic/gin"
	"net-share/pkg/bean"
	"net-share/server/service/client_host"
)

func (*api) List(c *gin.Context) {
	bean.Response.Success(c, "", client_host.Service.List())
}
