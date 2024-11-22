package client_host

import (
	"net-share/server/global"
	"net-share/server/service/registry"
)

func (*service) Delete(code string) error {
	host, err := global.ClientHostFs.Query(code)
	if err != nil {
		return err
	}
	if err := global.ClientHostFs.Delete(host.GetKey()); err != nil {
		return err
	}
	_ = global.ClientHostDomainFs.Delete(host.DomainPrefix)

	registry.ClientRegistry.Get(host.ClientCode).DelHost(host.Code)
	registry.UpdateIngress()
	return nil
}
