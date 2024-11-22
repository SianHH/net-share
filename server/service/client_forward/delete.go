package client_forward

import (
	"errors"
	"go.uber.org/zap"
	"net-share/server/global"
	"net-share/server/model"
	"net-share/server/service/ports"
	"net-share/server/service/registry"
)

func (*service) Delete(code string) error {
	var forward model.ClientForward

	if err := global.ClientForwardFs.Delete(code); err != nil {
		global.Logger.Error("删除ClientForward失败", zap.Error(err))
		return errors.New("删除失败")
	}
	ports.PushPort(forward.Port)
	registry.ClientRegistry.Get(forward.ClientCode).DelForward(forward.Code)
	return nil
}
