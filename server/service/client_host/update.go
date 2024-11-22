package client_host

import (
	"errors"
	"go.uber.org/zap"
	"net-share/pkg/utils"
	"net-share/server/constant"
	"net-share/server/global"
	"net-share/server/model"
	"net-share/server/service/registry"
	"regexp"
	"time"
)

type UpdateRequest struct {
	Code         string `binding:"required" json:"code"`
	Name         string `binding:"required" json:"name"`
	TargetIp     string `binding:"required" json:"targetIp"`
	TargetPort   string `binding:"required" json:"targetPort"`
	DomainPrefix string `binding:"required" json:"domainPrefix"`

	RateLimiter int `json:"rateLimiter"`
}

func (*service) Update(params UpdateRequest) (Item, error) {
	if !utils.ValidateLocalIP(params.TargetIp) || !utils.ValidatePort(params.TargetPort) {
		return Item{}, errors.New("IP或PORT不合法")
	}
	compile := regexp.MustCompile("^[a-z0-9]+$")
	if !compile.MatchString(params.DomainPrefix) {
		return Item{}, errors.New("域名前缀只能包含小写字母和数字")
	}

	host, err := global.ClientHostFs.Query(params.Code)
	if err != nil {
		return Item{}, err
	}

	client, err := global.ClientFs.Query(host.ClientCode)
	if err != nil {
		return Item{}, err
	}

	rateLimiter := params.RateLimiter

	oldDomainPrefix := host.DomainPrefix
	if oldDomainPrefix != params.DomainPrefix {
		if err = global.ClientHostDomainFs.Create(model.Base{
			Code:      params.DomainPrefix,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}); err != nil {
			return Item{}, err
		}
	}

	host.Name = params.Name
	host.Target = params.TargetIp + ":" + params.TargetPort
	host.DomainPrefix = params.DomainPrefix
	host.RateLimiter = rateLimiter
	host.AuthUser = utils.RandStr(12, utils.AllDict)
	host.AuthPwd = utils.RandStr(12, utils.AllDict)
	host.UpdatedAt = time.Now()

	if err := global.ClientHostFs.Update(host); err != nil {
		global.Logger.Error("修改ClientHost失败", zap.Error(err))
		_ = global.ClientHostDomainFs.Delete(params.DomainPrefix)
		return Item{}, errors.New("修改失败")
	}
	_ = global.ClientHostDomainFs.Delete(oldDomainPrefix)
	registry.ClientRegistry.Get(host.ClientCode).RunHost(host.Code, false)
	registry.UpdateIngress()
	registry.UpdateAuthers()

	ip, port := host.GetTargetIpAndPort()
	return Item{
		Code:           host.Code,
		Name:           host.Name,
		DomainPrefix:   host.DomainPrefix,
		DomainFull:     host.DomainPrefix + "." + global.App.Domain,
		BaseDomain:     "." + global.App.Domain,
		TargetIp:       ip,
		TargetPort:     port,
		ClientCode:     host.ClientCode,
		ClientName:     client.Name,
		ClientIsOnline: utils.TrinaryOperation(global.Cache.GetString(constant.CacheClientHeartbeatKey+host.ClientCode) == "online", 1, 2),
		RateLimiter:    host.RateLimiter,
	}, nil
}
