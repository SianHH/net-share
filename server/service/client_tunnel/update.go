package client_tunnel

import (
	"errors"
	"go.uber.org/zap"
	"net-share/pkg/utils"
	"net-share/server/constant"
	"net-share/server/global"
	"net-share/server/service/registry"
	"time"
)

type UpdateRequest struct {
	Code        string `binding:"required" json:"code"`
	Name        string `binding:"required" json:"name"`
	TargetIp    string `binding:"required" json:"targetIp"`
	TargetPort  string `binding:"required" json:"targetPort"`
	RateLimiter int    `json:"rateLimiter"`
}

func (*service) Update(params UpdateRequest) (Item, error) {
	if !utils.ValidateLocalIP(params.TargetIp) || !utils.ValidatePort(params.TargetPort) {
		return Item{}, errors.New("IP或PORT不合法")
	}
	tunnel, err := global.ClientTunnelFs.Query(params.Code)
	if err != nil {
		return Item{}, err
	}

	client, err := global.ClientFs.Query(tunnel.ClientCode)
	if err != nil {
		return Item{}, err
	}

	rateLimiter := params.RateLimiter

	tunnel.Name = params.Name
	tunnel.Target = params.TargetIp + ":" + params.TargetPort
	tunnel.RateLimiter = rateLimiter
	tunnel.UpdatedAt = time.Now()
	tunnel.AuthUser = utils.RandStr(12, utils.AllDict)
	tunnel.AuthPwd = utils.RandStr(12, utils.AllDict)
	if err := global.ClientTunnelFs.Update(tunnel); err != nil {
		global.Logger.Error("修改ClientTunnel失败", zap.Error(err))
		return Item{}, errors.New("修改失败")
	}
	registry.ClientRegistry.Get(tunnel.ClientCode).RunTunnel(tunnel.Code, false)
	registry.UpdateAuthers()
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
