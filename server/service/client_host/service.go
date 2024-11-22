package client_host

import (
	"net-share/server/framework/hook"
	"net-share/server/global"
	"time"
)

type service struct {
}

var Service = &service{}

func init() {
	hook.AddServerBeforeHookFunc(func() {
		go func() {
			for {
				var domainPrefixMap = make(map[string]bool)
				for _, host := range global.ClientHostFs.QueryAll() {
					domainPrefixMap[host.DomainPrefix] = true
				}
				var domainPrefixList []string
				for _, domainPrefix := range global.ClientHostDomainFs.QueryAll() {
					domainPrefixList = append(domainPrefixList, domainPrefix.Code)
				}
				for _, domainPrefix := range domainPrefixList {
					if !domainPrefixMap[domainPrefix] {
						_ = global.ClientHostDomainFs.Delete(domainPrefix)
					}
				}
				time.Sleep(time.Hour)
			}
		}()
	})
}
