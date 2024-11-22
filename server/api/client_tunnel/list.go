package client_tunnel

import (
	"github.com/gin-gonic/gin"
	"net-share/pkg/bean"
	"net-share/server/service/client_tunnel"
)

func (*api) List(c *gin.Context) {
	bean.Response.Success(c, "", client_tunnel.Service.List())
}
