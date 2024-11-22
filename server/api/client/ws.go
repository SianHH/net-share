package client

import (
	"github.com/gin-gonic/gin"
	"net-share/server/service/client"
)

func (*api) Ws(c *gin.Context) {
	client.Service.Ws(c)
}
