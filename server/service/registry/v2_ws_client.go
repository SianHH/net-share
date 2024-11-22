package registry

import (
	"encoding/json"
	"github.com/go-gost/x/config"
	"github.com/google/uuid"
	"github.com/lxzan/gws"
	"net-share/pkg/utils"
	"net-share/pkg/ws_msg_data"
	"net-share/server/constant"
	"net-share/server/global"
	"net-share/server/model"
	"time"
)

func NewV2Client(code string, svr *V2ClientWs) *V2Client {
	return &V2Client{
		code:      code,
		svr:       svr,
		isRunning: true,
	}
}

type V2Client struct {
	code      string
	svr       *V2ClientWs
	isRunning bool
}

func (c *V2Client) RunHost(code string, force bool) {
	if !c.checkExec() {
		return
	}
	type RunTunnelRequest struct {
		Code      string               `json:"code"`
		UpdatedAt int64                `json:"updatedAt"`
		Service   config.ServiceConfig `json:"service,omitempty"`
		Chain     config.ChainConfig   `json:"chain,omitempty"`
		Limiter   config.LimiterConfig `json:"limiter,omitempty"`
	}
	var host model.ClientHost
	host, err := global.ClientHostFs.Query(code)
	if err != nil {
		return
	}

	uid := uuid.NewString()
	updatedAt := time.Now().UnixNano()
	metadata := map[string]any{
		"tunnel.id": host.Code,
	}
	//for k, v := range constant.TunnelClientConfigMap {
	//	metadata[k] = v
	//}

	msgBody := ws_msg_data.NewMessageBody("runTunnel", uid, true, RunTunnelRequest{
		Code:      host.Code,
		UpdatedAt: updatedAt,
		Service: config.ServiceConfig{
			Name: host.Code,
			Addr: ":0",
			Handler: &config.HandlerConfig{Type: "rtcp", Metadata: map[string]any{
				"sniffing": true,
			}},
			Listener: &config.ListenerConfig{Type: "rtcp", Chain: host.Code},
			Forwarder: &config.ForwarderConfig{
				Nodes: []*config.ForwardNodeConfig{
					{
						Name: host.Target,
						Addr: host.Target,
						HTTP: &config.HTTPNodeConfig{
							Host: host.DomainPrefix + "." + global.App.Domain,
						},
					},
					{
						Name: host.Target,
						Addr: host.Target,
						HTTP: &config.HTTPNodeConfig{
							Host: host.DomainPrefix + "." + global.App.Domain,
							ResponseHeader: map[string]string{
								"Cache-Control": "private, max-age=86400, stale-while-revalidate=604800",
							},
						},
						Matcher: &config.NodeMatcherConfig{
							Rule: "!Header(`key`) && PathRegexp(`\\.(css|js|jpeg|jpg|png|gif|bmp|webp|svg|ttf|woff|woff2)$`)",
						},
					},
				},
			},
			Limiter: utils.TrinaryOperation(host.RateLimiter == 0, "", host.Code),
		},
		Chain: config.ChainConfig{
			Name: host.Code,
			Hops: []*config.HopConfig{
				{
					Nodes: []*config.NodeConfig{
						{
							Addr: ":" + global.App.HostPort,
							Connector: &config.ConnectorConfig{
								Type:     "tunnel",
								Metadata: metadata,
							},
							Dialer: &config.DialerConfig{
								Type: "tls",
							},
						},
					},
				},
			},
		},
		Limiter: config.LimiterConfig{
			Name:   host.Code,
			Limits: host.GetLimits(),
		},
	})
	c.WriterMessage(msgBody)
}

func (c *V2Client) DelHost(code string) {
	if !c.checkExec() {
		return
	}
	uid := uuid.NewString()
	msgBody := ws_msg_data.NewMessageBody("delTunnel", uid, true, map[string]string{"code": code})
	c.WriterMessage(msgBody)
}

func (c *V2Client) RunForward(code string, force bool) {
	if !c.checkExec() {
		return
	}
	type RunTunnelRequest struct {
		Code      string               `json:"code"`
		UpdatedAt int64                `json:"updatedAt"`
		Service   config.ServiceConfig `json:"service,omitempty"`
		Chain     config.ChainConfig   `json:"chain,omitempty"`
		Limiter   config.LimiterConfig `json:"limiter,omitempty"`
	}
	var forward model.ClientForward
	forward, err := global.ClientForwardFs.Query(code)
	if err != nil {
		return
	}
	uid := uuid.NewString()
	{
		tcpCode := forward.Code
		msgBody := ws_msg_data.NewMessageBody("runTunnel", uid, true, RunTunnelRequest{
			Code:      tcpCode,
			UpdatedAt: time.Now().UnixNano(),
			Service: config.ServiceConfig{
				Name: tcpCode,
				Addr: ":" + forward.Port,
				Handler: &config.HandlerConfig{
					Type: "rtcp",
					Metadata: map[string]any{
						"keepAlive": true,
					},
				},
				Listener: &config.ListenerConfig{
					Type:  "rtcp",
					Chain: tcpCode,
					Metadata: map[string]any{
						"keepAlive": true,
					},
				},
				Forwarder: &config.ForwarderConfig{
					Nodes: []*config.ForwardNodeConfig{
						{
							Name: forward.Target,
							Addr: forward.Target,
						},
					},
				},
				Limiter: utils.TrinaryOperation(forward.RateLimiter == 0, "", forward.Code),
			},
			Chain: config.ChainConfig{
				Name: tcpCode,
				Hops: []*config.HopConfig{
					{
						Nodes: []*config.NodeConfig{
							{
								Addr: global.App.Ip + ":" + global.App.ForwardPort,
								Connector: &config.ConnectorConfig{
									Type: "relay",
									Metadata: map[string]any{
										"nodelay": utils.TrinaryOperation(forward.Nodelay == 1, true, false),
									},
								},
								Dialer: &config.DialerConfig{
									Type: "tls",
								},
							},
						},
					},
				},
			},
			Limiter: config.LimiterConfig{
				Name:   tcpCode,
				Limits: forward.GetLimits(),
			},
		})
		c.WriterMessage(msgBody)
	}

	{
		udpCode := "udp_" + forward.Code
		msgBody := ws_msg_data.NewMessageBody("runTunnel", uid, true, RunTunnelRequest{
			Code:      udpCode,
			UpdatedAt: time.Now().UnixNano(),
			Service: config.ServiceConfig{
				Name: udpCode,
				Addr: ":" + forward.Port,
				Handler: &config.HandlerConfig{
					Type: "rudp",
					Metadata: map[string]any{
						"keepAlive": true,
					},
				},
				Listener: &config.ListenerConfig{
					Type:  "rudp",
					Chain: udpCode,
					Metadata: map[string]any{
						"keepAlive": true,
					},
				},
				Forwarder: &config.ForwarderConfig{
					Nodes: []*config.ForwardNodeConfig{
						{
							Name: forward.Target,
							Addr: forward.Target,
						},
					},
				},
				Limiter: utils.TrinaryOperation(forward.RateLimiter == 0, "", forward.Code),
			},
			Chain: config.ChainConfig{
				Name: udpCode,
				Hops: []*config.HopConfig{
					{
						Nodes: []*config.NodeConfig{
							{
								Addr: global.App.Ip + ":" + global.App.ForwardPort,
								Connector: &config.ConnectorConfig{
									Type: "relay",
									Metadata: map[string]any{
										"nodelay": utils.TrinaryOperation(forward.Nodelay == 1, true, false),
									},
								},
								Dialer: &config.DialerConfig{
									Type: "tls",
								},
							},
						},
					},
				},
			},
			Limiter: config.LimiterConfig{
				Name:   udpCode,
				Limits: forward.GetLimits(),
			},
		})
		c.WriterMessage(msgBody)
	}
}

func (c *V2Client) DelForward(code string) {
	if !c.checkExec() {
		return
	}
	uid := uuid.NewString()
	{
		msgBody := ws_msg_data.NewMessageBody("delTunnel", uid, true, map[string]string{"code": code})
		c.WriterMessage(msgBody)
	}
	{
		msgBody := ws_msg_data.NewMessageBody("delTunnel", uid, true, map[string]string{"code": "udp_" + code})
		c.WriterMessage(msgBody)
	}
}

func (c *V2Client) RunTunnel(code string, force bool) {
	if !c.checkExec() {
		return
	}
	type RunTunnelRequest struct {
		Code      string               `json:"code"`
		UpdatedAt int64                `json:"updatedAt"`
		Service   config.ServiceConfig `json:"service,omitempty"`
		Chain     config.ChainConfig   `json:"chain,omitempty"`
		Limiter   config.LimiterConfig `json:"limiter,omitempty"`
	}
	tunnel, err := global.ClientTunnelFs.Query(code)
	if err != nil {
		return
	}
	uid := uuid.NewString()
	metadata := map[string]any{
		"tunnel.id": tunnel.Key,
	}
	//for k, v := range constant.TunnelClientConfigMap {
	//	metadata[k] = v
	//}

	{
		tcpCode := tunnel.Code
		msgBody := ws_msg_data.NewMessageBody("runTunnel", uid, true, RunTunnelRequest{
			Code:      tcpCode,
			UpdatedAt: time.Now().UnixNano(),
			Service: config.ServiceConfig{
				Name: tcpCode,
				Addr: ":0",
				Handler: &config.HandlerConfig{
					Type: "rtcp",
				},
				Listener: &config.ListenerConfig{
					Type:  "rtcp",
					Chain: tcpCode,
				},
				Forwarder: &config.ForwarderConfig{
					Nodes: []*config.ForwardNodeConfig{
						{
							Name: tunnel.Target,
							Addr: tunnel.Target,
						},
					},
				},
				Limiter: utils.TrinaryOperation(tunnel.RateLimiter == 0, "", tunnel.Code),
			},
			Chain: config.ChainConfig{
				Name: tcpCode,
				Hops: []*config.HopConfig{
					{
						Nodes: []*config.NodeConfig{
							{
								Addr: global.App.Ip + ":" + global.App.HostPort,
								Connector: &config.ConnectorConfig{
									Type:     "tunnel",
									Metadata: metadata,
								},
								Dialer: &config.DialerConfig{
									Type: "tls",
								},
							},
						},
					},
				},
			},
			Limiter: config.LimiterConfig{
				Name:   tcpCode,
				Limits: tunnel.GetLimits(),
			},
		})
		c.WriterMessage(msgBody)
	}

	{
		udpCode := "udp_" + tunnel.Code
		msgBody := ws_msg_data.NewMessageBody("runTunnel", uid, true, RunTunnelRequest{
			Code:      udpCode,
			UpdatedAt: time.Now().UnixNano(),
			Service: config.ServiceConfig{
				Name: udpCode,
				Addr: ":0",
				Handler: &config.HandlerConfig{
					Type: "rudp",
				},
				Listener: &config.ListenerConfig{
					Type:  "rudp",
					Chain: udpCode,
				},
				Forwarder: &config.ForwarderConfig{
					Nodes: []*config.ForwardNodeConfig{
						{
							Name: tunnel.Target,
							Addr: tunnel.Target,
						},
					},
				},
				Limiter: utils.TrinaryOperation(tunnel.RateLimiter == 0, "", tunnel.Code),
			},
			Chain: config.ChainConfig{
				Name: udpCode,
				Hops: []*config.HopConfig{
					{
						Nodes: []*config.NodeConfig{
							{
								Addr: global.App.Ip + ":" + global.App.HostPort,
								Connector: &config.ConnectorConfig{
									Type:     "tunnel",
									Metadata: metadata,
								},
								Dialer: &config.DialerConfig{
									Type: "tls",
								},
							},
						},
					},
				},
			},
			Limiter: config.LimiterConfig{
				Name:   udpCode,
				Limits: tunnel.GetLimits(),
			},
		})
		c.WriterMessage(msgBody)
	}
}

func (c *V2Client) DelTunnel(code string) {
	if !c.checkExec() {
		return
	}
	uid := uuid.NewString()
	{
		msgBody := ws_msg_data.NewMessageBody("delTunnel", uid, true, map[string]string{"code": code})
		c.WriterMessage(msgBody)
	}
	{
		msgBody := ws_msg_data.NewMessageBody("delTunnel", uid, true, map[string]string{"code": "udp_" + code})
		c.WriterMessage(msgBody)
	}
}

func (c *V2Client) checkExec() bool {
	if c.svr == nil || !c.isRunning {
		return false
	}
	return true
}

func (c *V2Client) Stop(reason string) {
	if !c.checkExec() {
		return
	}
	uid := uuid.NewString()
	msgBody := ws_msg_data.MessageBody{
		Callback: "stop",
		UUID:     uid,
		Aes:      false,
	}
	_ = msgBody.Marshal(map[string]string{
		"msg": reason,
	})
	c.WriterMessage(msgBody)
	_ = c.svr.conn.WriteClose(1000, nil)
}

func (c *V2Client) Init() {
	if !c.checkExec() {
		return
	}
	uid := uuid.NewString()
	msgBody := ws_msg_data.NewMessageBody("init", uid, false, map[string]string{})
	c.WriterMessage(msgBody)
}

func (c *V2Client) WriterMessage(body ws_msg_data.MessageBody) {
	if c.svr == nil || c.svr.conn == nil || !c.svr.isRunning {
		return
	}
	marshal, _ := json.Marshal(body)
	_ = c.svr.conn.SetDeadline(time.Now().Add(time.Second * 15))
	_ = c.svr.conn.WriteString(string(marshal))
}

func NewV2ClientWs(code string) *V2ClientWs {
	return &V2ClientWs{
		code:      code,
		isRunning: true,
	}
}

type V2ClientWs struct {
	code      string
	conn      *gws.Conn
	isRunning bool
}

func (c *V2ClientWs) OnOpen(socket *gws.Conn) {
	c.conn = socket
	c.isRunning = true
	_ = socket.SetDeadline(time.Now().Add(time.Second * 15))
	global.Cache.SetString(constant.CacheClientHeartbeatKey+c.code, "online", time.Second*10)
}

func (c *V2ClientWs) OnClose(socket *gws.Conn, err error) {
	c.isRunning = false
	global.Cache.Delete(constant.CacheClientHeartbeatKey + c.code)
}

func (c *V2ClientWs) OnPing(socket *gws.Conn, payload []byte) {
	_ = socket.SetDeadline(time.Now().Add(time.Second * 15))
	_ = socket.WritePong(nil)
	global.Cache.SetString(constant.CacheClientHeartbeatKey+c.code, "online", time.Second*10)
}

func (c *V2ClientWs) OnPong(socket *gws.Conn, payload []byte) {}

func (c *V2ClientWs) OnMessage(socket *gws.Conn, message *gws.Message) {
	defer message.Close()
	global.Cache.SetString(constant.CacheClientHeartbeatKey+c.code, "online", time.Second*10)
	var data ws_msg_data.MessageBody
	_ = json.Unmarshal(message.Bytes(), &data)
	switch data.Callback {
	case "init":

	}
}
