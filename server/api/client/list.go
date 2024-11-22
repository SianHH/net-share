package client

import (
	"github.com/gin-gonic/gin"
	"net-share/pkg/bean"
	"net-share/server/service/client"
)

func (*api) List(c *gin.Context) {
	bean.Response.Success(c, "", client.Service.List())
}
