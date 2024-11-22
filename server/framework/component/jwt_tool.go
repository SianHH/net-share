package component

import (
	"net-share/pkg/jwt"
	"net-share/server/global"
)

func InitJwtTool() {
	global.JwtTool = jwt.NewTool(global.App.JwtKey)
}
