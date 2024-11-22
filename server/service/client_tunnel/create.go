package client_tunnel

import (
	"errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net-share/pkg/utils"
	"net-share/server/constant"
	"net-share/server/global"
	"net-share/server/model"
	"net-share/server/service/registry"
	"time"
)

type CreateRequest struct {
	Name       string `binding:"required" json:"name"`
	TargetIp   string `binding:"required" json:"targetIp"`
	TargetPort string `binding:"required" json:"targetPort"`
	ClientCode string `binding:"required" json:"clientCode"`

	Nodelay     int `json:"nodelay"`
	RateLimiter int `json:"rateLimiter"`
}

func (*service) Create(params CreateRequest) (Item, error) {
	if !utils.ValidateLocalIP(params.TargetIp) || !utils.ValidatePort(params.TargetPort) {
		return Item{}, errors.New("IP或PORT不合法")
	}

	var client model.Client
	client, err := global.ClientFs.Query(params.ClientCode)
	if err != nil {
		return Item{}, errors.New("客户端不存在")
	}

	rateLimiter := params.RateLimiter

	tunnel := model.ClientTunnel{
		Base: model.Base{
			Code:      uuid.NewString(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name:        params.Name,
		Target:      params.TargetIp + ":" + params.TargetPort,
		ClientCode:  params.ClientCode,
		Key:         uuid.NewString(),
		RateLimiter: rateLimiter,
		AuthUser:    utils.RandStr(12, utils.AllDict),
		AuthPwd:     utils.RandStr(12, utils.AllDict),
	}
	if err := global.ClientTunnelFs.Create(tunnel); err != nil {
		global.Logger.Error("新增ClientTunnel失败", zap.Error(err))
		return Item{}, err
	}
	registry.ClientRegistry.Get(tunnel.ClientCode).RunTunnel(tunnel.Code, false)
	registry.UpdateIngress()
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
