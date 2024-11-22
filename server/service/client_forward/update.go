package client_forward

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
	Code       string `binding:"required" json:"code"`
	Name       string `binding:"required" json:"name"`
	TargetIp   string `binding:"required" json:"targetIp"`
	TargetPort string `binding:"required" json:"targetPort"`

	Nodelay     int `json:"nodelay"`
	RateLimiter int `json:"rateLimiter"`
}

func (*service) Update(params UpdateRequest) (Item, error) {
	if !utils.ValidateLocalIP(params.TargetIp) || !utils.ValidatePort(params.TargetPort) {
		return Item{}, errors.New("IP或PORT不合法")
	}
	forward, err := global.ClientForwardFs.Query(params.Code)
	if err != nil {
		return Item{}, err
	}

	client, err := global.ClientFs.Query(forward.ClientCode)
	if err != nil {
		return Item{}, err
	}

	rateLimiter := params.RateLimiter

	forward.Name = params.Name
	forward.Target = params.TargetIp + ":" + params.TargetPort
	forward.Nodelay = params.Nodelay
	forward.RateLimiter = rateLimiter
	forward.UpdatedAt = time.Now()
	forward.AuthUser = utils.RandStr(12, utils.AllDict)
	forward.AuthPwd = utils.RandStr(12, utils.AllDict)
	if err := global.ClientForwardFs.Update(forward); err != nil {
		global.Logger.Error("修改ClientForward失败", zap.Error(err))
		return Item{}, errors.New("修改失败")
	}
	registry.ClientRegistry.Get(forward.ClientCode).RunForward(forward.Code, false)
	registry.UpdateAuthers()
	ip, port := forward.GetTargetIpAndPort()
	return Item{
		Code:           forward.Code,
		Name:           forward.Name,
		TargetIp:       ip,
		TargetPort:     port,
		Port:           forward.Port,
		Ip:             global.App.Ip,
		ClientCode:     forward.ClientCode,
		ClientName:     client.Name,
		ClientIsOnline: utils.TrinaryOperation(global.Cache.GetString(constant.CacheClientHeartbeatKey+forward.ClientCode) == "online", 1, 2),
		Nodelay:        forward.Nodelay,
		RateLimiter:    forward.RateLimiter,
	}, nil
}
