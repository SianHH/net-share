package component

import (
	"go.uber.org/zap"
	"net-share/pkg/file_storage"
	"net-share/server/global"
	"net-share/server/model"
	"os"
)

func InitRepo() {
	var err error
	global.ClientFs, err = file_storage.NewFileStorage(model.Client{}, "data/client.json", 10)
	if err != nil {
		global.Logger.Error("ClientFs初始化失败", zap.Error(err))
		os.Exit(1)
	}

	global.ClientForwardFs, err = file_storage.NewFileStorage(model.ClientForward{}, "data/client_forward.json", 10)
	if err != nil {
		global.Logger.Error("ClientForwardFs初始化失败", zap.Error(err))
		os.Exit(1)
	}

	global.ClientHostFs, err = file_storage.NewFileStorage(model.ClientHost{}, "data/client_host.json", 10)
	if err != nil {
		global.Logger.Error("ClientHostFs初始化失败", zap.Error(err))
		os.Exit(1)
	}

	global.ClientHostDomainFs, err = file_storage.NewFileStorage(model.Base{}, "data/client_host_domain.json", 10)
	if err != nil {
		global.Logger.Error("ClientHostDomainFs初始化失败", zap.Error(err))
		os.Exit(1)
	}
}
