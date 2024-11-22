package client

import (
	"net-share/pkg/utils"
	"net-share/server/constant"
	"net-share/server/global"
)

type ListItem struct {
	Code         string `json:"code"`
	Name         string `json:"name"`
	Key          string `json:"key"`
	ForwardTotal int    `json:"forwardTotal"`
	HostTotal    int    `json:"hostTotal"`
	IsOnline     int    `json:"isOnline"`
}

func (*service) List() (list []ListItem) {
	for _, client := range global.ClientFs.QueryAll() {
		list = append(list, ListItem{
			Code: client.Code,
			Name: client.Name,
			Key:  client.Key,
			//ForwardTotal: len(client.ForwardList),
			//HostTotal:    len(client.HostList),
			IsOnline: utils.TrinaryOperation(global.Cache.GetString(constant.CacheClientHeartbeatKey+client.Code) == "online", 1, 2),
		})
	}
	return list
}
