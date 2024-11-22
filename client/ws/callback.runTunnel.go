package ws

import (
	"github.com/go-gost/core/logger"
	"github.com/go-gost/x/config"
	chainParser "github.com/go-gost/x/config/parsing/chain"
	limiterParser "github.com/go-gost/x/config/parsing/limiter"
	serviceParser "github.com/go-gost/x/config/parsing/service"
	"github.com/go-gost/x/registry"
	"net-share/pkg/ws_msg_data"
)

type RunTunnelRequest struct {
	Code      string               `json:"code"`
	UpdatedAt int64                `json:"updatedAt"`
	Service   config.ServiceConfig `json:"service,omitempty"`
	Chain     config.ChainConfig   `json:"chain,omitempty"`
	Limiter   config.LimiterConfig `json:"limiter,omitempty"`
}

func RunTunnel(data ws_msg_data.MessageBody) *ws_msg_data.MessageBody {
	var param RunTunnelRequest
	_ = data.Unmarshal(&param)
	var errs []string
	if !CheckUpdate(param.Code, param.UpdatedAt) {
		resp := ws_msg_data.NewMessageBody(data.Callback, data.UUID, false, errs)
		return &resp
	}
	Store(param.Code, param.UpdatedAt)
	//
	old := registry.ServiceRegistry().Get(param.Service.Name)
	if old != nil {
		_ = old.Close()
	}
	registry.ServiceRegistry().Unregister(param.Service.Name)
	service, err := serviceParser.ParseService(&param.Service)
	if err != nil {
		errs = append(errs, err.Error())
	} else {
		go service.Serve()
		_ = registry.ServiceRegistry().Register(param.Service.Name, service)
	}

	//
	chain, err := chainParser.ParseChain(&param.Chain, logger.Default())
	if err != nil {
		errs = append(errs, err.Error())
	}
	registry.ChainRegistry().Unregister(param.Chain.Name)
	_ = registry.ChainRegistry().Register(param.Chain.Name, chain)

	//
	registry.TrafficLimiterRegistry().Unregister(param.Limiter.Name)
	parseRateLimiter := limiterParser.ParseTrafficLimiter(&param.Limiter)
	_ = registry.TrafficLimiterRegistry().Register(param.Limiter.Name, parseRateLimiter)

	resp := ws_msg_data.NewMessageBody(data.Callback, data.UUID, false, errs)
	return &resp
}
