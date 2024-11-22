package client_forward

import (
	"errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net-share/pkg/utils"
	"net-share/server/constant"
	"net-share/server/global"
	"net-share/server/model"
	"net-share/server/service/ports"
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

	port, err := ports.PullPort()
	if err != nil {
		return Item{}, err
	}
	forward := model.ClientForward{
		Base: model.Base{
			Code:      uuid.NewString(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name:        params.Name,
		Target:      params.TargetIp + ":" + params.TargetPort,
		Port:        port,
		Nodelay:     params.Nodelay,
		ClientCode:  params.ClientCode,
		RateLimiter: rateLimiter,
	}
	if err := global.ClientForwardFs.Create(forward); err != nil {
		ports.PushPort(port)
		global.Logger.Error("新增ClientForward失败", zap.Error(err))
		return Item{}, err
	}
	registry.ClientRegistry.Get(forward.ClientCode).RunForward(forward.Code, false)

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
