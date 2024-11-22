package client

import (
	"errors"
	"go.uber.org/zap"
	"net-share/server/global"
	"net-share/server/service/registry"
)

func (*service) Delete(code string) error {
	if err := global.ClientFs.Delete(code); err != nil {
		global.Logger.Error("删除client失败", zap.Error(err))
		return errors.New("删除失败")
	}

	// 清理占用的资源
	var hostCodeList []string
	var hostDomainPrefixList []string
	for _, host := range global.ClientHostFs.QueryAll() {
		if host.ClientCode == code {
			hostCodeList = append(hostCodeList, host.Code)
			hostDomainPrefixList = append(hostDomainPrefixList, host.DomainPrefix)
		}
	}
	for _, hostCode := range hostCodeList {
		global.ClientHostFs.Delete(hostCode)
	}

	for _, domainPrefix := range hostDomainPrefixList {
		global.ClientHostDomainFs.Delete(domainPrefix)
	}

	var forwardCodeList []string
	for _, forward := range global.ClientForwardFs.QueryAll() {
		if forward.ClientCode == code {
			forwardCodeList = append(forwardCodeList, forward.Code)
		}
	}
	for _, forwardCode := range forwardCodeList {
		global.ClientForwardFs.Delete(forwardCode)
	}

	registry.ClientRegistry.Get(code).Stop("客户端已被删除，停止进程")
	registry.UpdateIngress()
	return nil
}
