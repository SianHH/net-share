package router

import (
	"github.com/gin-gonic/gin"
)

func SetRouter(engine *gin.Engine) {
	setAuthRouter(engine)
	setClientRouter(engine)
	setClientHostRouter(engine)
	setClientForwardRouter(engine)
	setClientTunnelRouter(engine)
}
