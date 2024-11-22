package ws

import (
	"github.com/go-gost/x/registry"
	"net-share/pkg/ws_msg_data"
)

type DelTunnelRequest struct {
	Code string `json:"code"`
}

func DelTunnel(data ws_msg_data.MessageBody) *ws_msg_data.MessageBody {
	var param DelTunnelRequest
	_ = data.Unmarshal(&param)
	Remove(param.Code)

	service := registry.ServiceRegistry().Get(param.Code)
	if service != nil {
		_ = service.Close()
	}
	registry.ServiceRegistry().Unregister(param.Code)

	registry.ChainRegistry().Unregister(param.Code)
	registry.TrafficLimiterRegistry().Unregister(param.Code)

	resp := ws_msg_data.NewMessageBody(data.Callback, data.UUID, false, "success")
	return &resp
}
