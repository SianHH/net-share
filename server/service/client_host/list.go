package client_host

import (
	"net-share/pkg/utils"
	"net-share/server/constant"
	"net-share/server/global"
)

type Item struct {
	Code           string `json:"code"`
	Name           string `json:"name"`
	DomainPrefix   string `json:"domainPrefix"`
	DomainFull     string `json:"domainFull"`
	BaseDomain     string `json:"baseDomain"`
	TargetIp       string `json:"targetIp"`
	TargetPort     string `json:"targetPort"`
	ClientCode     string `json:"clientCode"`
	ClientName     string `json:"clientName"`
	ClientIsOnline int    `json:"clientIsOnline"`

	RateLimiter int `json:"rateLimiter"`
}

func (*service) List() (list []Item) {
	for _, tunnel := range global.ClientHostFs.QueryAll() {
		client, _ := global.ClientFs.Query(tunnel.ClientCode)
		ip, port := tunnel.GetTargetIpAndPort()
		list = append(list, Item{
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
		})
	}
	return list
}
