package ws

import (
	"encoding/json"
	"fmt"
	"github.com/lxzan/gws"
	"net-share/pkg/ws_msg_data"
	"time"
)

const clientVersion = "v20241117"

type WsServer struct {
	ws  string
	key string
}

func (s *WsServer) OnOpen(socket *gws.Conn) {}

func (s *WsServer) OnClose(socket *gws.Conn, err error) {
	fmt.Println("conn close", err)
}

func (s *WsServer) OnPing(socket *gws.Conn, payload []byte) {}

func (s *WsServer) OnPong(socket *gws.Conn, payload []byte) {
	time.Sleep(time.Second * 5)
	_ = socket.WritePing(nil)
	fmt.Println("pong...")
}

func (s *WsServer) OnMessage(socket *gws.Conn, message *gws.Message) {
	var data ws_msg_data.MessageBody
	_ = json.Unmarshal(message.Bytes(), &data)
	var result *ws_msg_data.MessageBody
	switch data.Callback {
	case "stop":
		result = Stop(data)
	case "init":
		result = Init(data)
	case "runServer":
		result = RunServer(data)
	case "runIngress":
		result = RunIngress(data)
	case "runTunnel":
		result = RunTunnel(data)
	case "delTunnel":
		result = DelTunnel(data)
	}
	if result != nil {
		marshal, _ := json.Marshal(result)
		_ = socket.WriteMessage(gws.OpcodeText, marshal)
	}
}

func NewService(ws, key string) *WsServer {
	return &WsServer{
		ws:  ws,
		key: key,
	}
}
