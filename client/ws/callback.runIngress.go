package ws

import (
	"github.com/go-gost/x/config"
	ingressParser "github.com/go-gost/x/config/parsing/ingress"
	"github.com/go-gost/x/registry"
	"net-share/pkg/ws_msg_data"
)

type RunIngressRequest struct {
	Code      string               `json:"code"`
	UpdatedAt int64                `json:"updatedAt"`
	Ingress   config.IngressConfig `json:"ingress"`
}

func RunIngress(data ws_msg_data.MessageBody) *ws_msg_data.MessageBody {
	var param RunIngressRequest
	_ = data.Unmarshal(&param)
	var errs []string

	if !CheckUpdate(param.Code, param.UpdatedAt) {
		resp := ws_msg_data.NewMessageBody("runIngress", data.UUID, false, errs)
		return &resp
	}
	Store(param.Code, param.UpdatedAt)

	ingress := ingressParser.ParseIngress(&param.Ingress)
	registry.IngressRegistry().Unregister(param.Ingress.Name)
	_ = registry.IngressRegistry().Register(param.Ingress.Name, ingress)

	_ = config.OnUpdate(func(c *config.Config) error {
		for i := range c.Ingresses {
			if c.Ingresses[i].Name == param.Ingress.Name {
				c.Ingresses[i] = &param.Ingress
				break
			}
		}
		return nil
	})

	resp := ws_msg_data.NewMessageBody(data.Callback, data.UUID, false, errs)
	return &resp
}
