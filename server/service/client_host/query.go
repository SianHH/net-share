package client_host

import (
	"net-share/pkg/utils"
	"net-share/server/constant"
	"net-share/server/global"
)

func (*service) Query(code string) (Item, error) {
	tunnel, err := global.ClientHostFs.Query(code)
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
		DomainPrefix:   tunnel.DomainPrefix,
		DomainFull:     tunnel.DomainPrefix + "." + global.App.Domain,
		BaseDomain:     "." + global.App.Domain,
		TargetIp:       ip,
		TargetPort:     port,
		ClientCode:     tunnel.ClientCode,
		ClientName:     client.Name,
		ClientIsOnline: utils.TrinaryOperation(global.Cache.GetString(constant.CacheClientHeartbeatKey+tunnel.ClientCode) == "online", 1, 2),
		RateLimiter:    tunnel.RateLimiter,
	}, nil

}
