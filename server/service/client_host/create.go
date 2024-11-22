package client_host

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
	Name        string `binding:"required" json:"name"`
	TargetIp    string `binding:"required" json:"targetIp"`
	TargetPort  string `binding:"required" json:"targetPort"`
	ClientCode  string `binding:"required" json:"clientCode"`
	RateLimiter int    `json:"rateLimiter"`
}

func (*service) Create(params CreateRequest) (Item, error) {
	if !utils.ValidateLocalIP(params.TargetIp) || !utils.ValidatePort(params.TargetPort) {
		return Item{}, errors.New("IP或PORT不合法")
	}

	client, err := global.ClientFs.Query(params.ClientCode)
	if err != nil {
		return Item{}, err
	}

	rateLimiter := params.RateLimiter

	domainPrefix := utils.RandStr(12, utils.AllDict)
	if err := global.ClientHostDomainFs.Create(model.Base{
		Code:      domainPrefix,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}); err != nil {
		return Item{}, err
	}

	host := model.ClientHost{
		Base: model.Base{
			Code:      uuid.NewString(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name:         params.Name,
		Target:       params.TargetIp + ":" + params.TargetPort,
		DomainPrefix: domainPrefix,
		ClientCode:   params.ClientCode,
		RateLimiter:  rateLimiter,
		AuthUser:     utils.RandStr(12, utils.AllDict),
		AuthPwd:      utils.RandStr(12, utils.AllDict),
	}
	if err := global.ClientHostFs.Create(host); err != nil {
		global.ClientHostDomainFs.Delete(domainPrefix)
		global.Logger.Error("新增ClientHost失败", zap.Error(err))
		return Item{}, err
	}

	registry.ClientRegistry.Get(host.ClientCode).RunHost(host.Code, false)
	registry.UpdateIngress()
	registry.UpdateAuthers()
	return Item{
		Code:           host.Code,
		Name:           host.Name,
		DomainPrefix:   host.DomainPrefix,
		DomainFull:     host.DomainPrefix + "." + global.App.Domain,
		BaseDomain:     "." + global.App.Domain,
		TargetIp:       params.TargetIp,
		TargetPort:     params.TargetPort,
		ClientCode:     host.ClientCode,
		ClientName:     client.Name,
		ClientIsOnline: utils.TrinaryOperation(global.Cache.GetString(constant.CacheClientHeartbeatKey+host.ClientCode) == "online", 1, 2),
		RateLimiter:    host.RateLimiter,
	}, nil
}
