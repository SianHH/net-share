package client_tunnel

import (
	"errors"
	"go.uber.org/zap"
	"net-share/server/global"
	"net-share/server/model"
	"net-share/server/service/registry"
)

func (*service) Delete(code string) error {
	var tunnel model.ClientTunnel

	if err := global.ClientTunnelFs.Delete(code); err != nil {
		global.Logger.Error("删除ClientTunnel失败", zap.Error(err))
		return errors.New("删除失败")
	}
	registry.ClientRegistry.Get(tunnel.ClientCode).DelTunnel(tunnel.Code)
	return nil
}
