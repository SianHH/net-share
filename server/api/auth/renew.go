package auth

import (
	"github.com/gin-gonic/gin"
	"net-share/pkg/bean"
	"net-share/server/service/auth"
)

func (*api) Renew(c *gin.Context) {
	token, err := auth.Service.Renew(c)
	if err != nil {
		bean.Response.Fail(c, err.Error())
		return
	}
	bean.Response.Success(c, "", token)
}
