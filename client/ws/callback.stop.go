package ws

import (
	"fmt"
	"github.com/go-gost/x/registry"
	"net-share/pkg/ws_msg_data"
	"os"
	"time"
)

type StopRequest struct {
	Msg string `json:"msg"`
}

func Stop(data ws_msg_data.MessageBody) *ws_msg_data.MessageBody {
	for name, _ := range registry.ServiceRegistry().GetAll() {
		registry.ServiceRegistry().Unregister(name)
	}
	var params StopRequest
	_ = data.Unmarshal(&params)
	go func() {
		fmt.Println("3秒后，结束进程")
		fmt.Println("结束原因：", params.Msg)
		time.Sleep(time.Second * 3)
		//s.client.Stop()
		os.Exit(0)
	}()
	return &ws_msg_data.MessageBody{
		Callback: data.Callback,
		UUID:     data.UUID,
		Aes:      false,
		Data:     "success",
	}
}
