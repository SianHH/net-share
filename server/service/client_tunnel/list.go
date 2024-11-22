package client_tunnel

import (
	"net-share/pkg/utils"
	"net-share/server/constant"
	"net-share/server/global"
)

type Item struct {
	Code           string `json:"code"`
	Name           string `json:"name"`
	TargetIp       string `json:"targetIp"`
	TargetPort     string `json:"targetPort"`
	VKey           string `json:"vKey"`
	Ip             string `json:"ip"`
	ClientCode     string `json:"clientCode"`
	ClientName     string `json:"clientName"`
	ClientIsOnline int    `json:"clientIsOnline"`

	RateLimiter int `json:"rateLimiter"`
}

func (*service) List() (list []Item) {
	for _, tunnel := range global.ClientTunnelFs.QueryAll() {
		client, _ := global.ClientFs.Query(tunnel.ClientCode)
		ip, port := tunnel.GetTargetIpAndPort()
		list = append(list, Item{
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
		})
	}
	return list
}
