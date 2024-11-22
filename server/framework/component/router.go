package component

import (
	"github.com/gin-gonic/gin"
	"net-share/pkg/gin_zip_static"
	"net-share/server/frontend"
	"net-share/server/global"
	"net-share/server/router"
	"net-share/server/router/middleware"
	"net/http"
	"strings"
)

var engine *gin.Engine

func InitRouter() {
	if global.App.Mode == "dev" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	engine = gin.New()
	engine.Use(
		middleware.Cors(),
		middleware.NoCache(),
		middleware.Exception(),
	)
	router.SetRouter(engine)

	// 注册静态资源
	gin_zip_static.RegisterStaticFile(frontend.GetHtmlZipFile(), func(fileMap map[string][]byte) {
		for k, data := range fileMap {
			// 规避forRange复用k,data
			fileKey := k
			fileBytes := data
			ginStaticFilePath := strings.Replace(fileKey, "dist/", "", 1)
			if ginStaticFilePath == "" {
				continue
			}

			engine.GET(ginStaticFilePath, middleware.Cache(), func(c *gin.Context) {
				c.Data(http.StatusOK, gin_zip_static.MatchFile(fileKey), fileBytes)
			})
		}
		engine.NoRoute(func(c *gin.Context) {
			c.Data(http.StatusOK, "text/html; charset=utf-8", fileMap["dist/index.html"])
		})
	})

	engine.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": 404,
			"msg":  "资源不存在",
		})
	})
}
