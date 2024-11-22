package client_forward

import (
	"github.com/gin-gonic/gin"
	"net-share/pkg/bean"
	"net-share/server/service/client_forward"
)

func (*api) List(c *gin.Context) {
	bean.Response.Success(c, "", client_forward.Service.List())
}
