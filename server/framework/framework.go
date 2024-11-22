package framework

import (
	"fmt"
	"net-share/server/framework/component"
)

// [SOURCE] https://patorjk.com/software/taag/#p=display&h=0&v=0&f=ANSI%20Shadow&t=NET-SHARE
func init() {
	fmt.Println(`
███╗   ██╗███████╗████████╗      ███████╗██╗  ██╗ █████╗ ██████╗ ███████╗
████╗  ██║██╔════╝╚══██╔══╝      ██╔════╝██║  ██║██╔══██╗██╔══██╗██╔════╝
██╔██╗ ██║█████╗     ██║   █████╗███████╗███████║███████║██████╔╝█████╗  
██║╚██╗██║██╔══╝     ██║   ╚════╝╚════██║██╔══██║██╔══██║██╔══██╗██╔══╝  
██║ ╚████║███████╗   ██║         ███████║██║  ██║██║  ██║██║  ██║███████╗
╚═╝  ╚═══╝╚══════╝   ╚═╝         ╚══════╝╚═╝  ╚═╝╚═╝  ╚═╝╚═╝  ╚═╝╚══════╝
`)
}

func Run() {
	component.InitConfig()
	component.InitLogger()
	component.InitJwtTool()
	component.InitCache()
	component.InitRepo()
	component.InitRouter()
	component.InitServer()
}
