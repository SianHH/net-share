package ws

import (
	"fmt"
	"net-share/pkg/ws_msg_data"
)

type InitRequest struct {
	Version string `json:"version"` // 客户端版本号
}

func Init(data ws_msg_data.MessageBody) *ws_msg_data.MessageBody {
	var params InitRequest
	_ = data.Unmarshal(&params)
	fmt.Println("CONNECT SUCCESS")
	fmt.Println("Client Version:", clientVersion)
	fmt.Println("Server Version:", params.Version)
	result := &ws_msg_data.MessageBody{
		Callback: data.Callback,
		UUID:     data.UUID,
		Aes:      false,
	}
	_ = result.Marshal(map[string]string{})
	return result
}
