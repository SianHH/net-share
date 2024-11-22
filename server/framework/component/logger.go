package component

import (
	"net-share/pkg/logger"
	"net-share/server/global"
)

func InitLogger() {
	global.Logger = logger.NewLogger(global.App.Logger.File, global.App.Logger.Level, global.App.Logger.Console)
}
