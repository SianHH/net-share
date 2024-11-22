package ws

import (
	"github.com/go-gost/x/config"
	recorderParser "github.com/go-gost/x/config/parsing/recorder"
	serviceParser "github.com/go-gost/x/config/parsing/service"
	"github.com/go-gost/x/registry"
	"net-share/pkg/ws_msg_data"
)

type RunServerRequest struct {
	Code        string                 `json:"code"`
	UpdatedAt   int64                  `json:"updatedAt"`
	ServiceList []config.ServiceConfig `json:"serviceList"`
	Recorder    config.RecorderConfig  `json:"recorder"`
}

func RunServer(data ws_msg_data.MessageBody) *ws_msg_data.MessageBody {
	var param RunServerRequest
	_ = data.Unmarshal(&param)
	var errs []string

	if !CheckUpdate(param.Code, param.UpdatedAt) {
		resp := ws_msg_data.NewMessageBody(data.Callback, data.UUID, false, errs)
		return &resp
	}
	Store(param.Code, param.UpdatedAt)

	for _, service := range param.ServiceList {
		old := registry.ServiceRegistry().Get(service.Name)
		if old != nil {
			_ = old.Close()
		}
		registry.ServiceRegistry().Unregister(service.Name)

		parseService, err := serviceParser.ParseService(&service)
		if err != nil {
			errs = append(errs, err.Error())
			continue
		}
		go parseService.Serve()
		_ = registry.ServiceRegistry().Register(service.Name, parseService)
	}

	recorder := recorderParser.ParseRecorder(&param.Recorder)
	registry.RecorderRegistry().Unregister(param.Recorder.Name)
	_ = registry.RecorderRegistry().Register(param.Recorder.Name, recorder)

	resp := ws_msg_data.NewMessageBody(data.Callback, data.UUID, false, errs)
	return &resp
}
