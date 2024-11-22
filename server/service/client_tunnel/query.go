package client_tunnel

import (
	"net-share/pkg/utils"
	"net-share/server/constant"
	"net-share/server/global"
)

func (*service) Query(code string) (Item, error) {
	tunnel, err := global.ClientTunnelFs.Query(code)
	if err != nil {
		return Item{}, err
	}
	client, err := global.ClientFs.Query(tunnel.ClientCode)
	if err != nil {
		return Item{}, err
	}
	ip, port := tunnel.GetTargetIpAndPort()
	return Item{
		Code:           tunnel.Code,
		Name:           tunnel.Name,
		TargetIp:       ip,
		TargetPort:     port,
		VKey:           tunnel.Key,
		Ip:             global.App.Ip,
		ClientCode:     tunnel.ClientCode,
		ClientName:     client.Name,
		ClientIsOnline: utils.TrinaryOperation(global.Cache.GetString(constant.CacheClientHeartbeatKey+tunnel.ClientCode) == "online", 1, 2),
		RateLimiter:    tunnel.RateLimiter,
	}, nil
}
