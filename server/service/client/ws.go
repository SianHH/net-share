package client

import (
	"github.com/gin-gonic/gin"
	"github.com/lxzan/gws"
	"go.uber.org/zap"
	"net-share/server/constant"
	"net-share/server/global"
	"net-share/server/model"
	"net-share/server/service/registry"
	"time"
)

const (
	PingInterval = 5 * time.Second
	PingWait     = 10 * time.Second
)

func (*service) Ws(c *gin.Context) {
	key := c.GetHeader("key")
	client, err := global.ClientFs.QueryFilter(func(client model.Client) bool {
		return client.Key == key
	})
	if err != nil {
		return
	}

	handler := registry.NewV2ClientWs(client.Code)
	upgrader := gws.NewUpgrader(handler, &gws.ServerOption{
		ParallelEnabled:   true,                                 // 开启并行消息处理
		Recovery:          gws.Recovery,                         // 开启异常恢复
		PermessageDeflate: gws.PermessageDeflate{Enabled: true}, // 开启压缩
	})
	upgrade, err := upgrader.Upgrade(c.Writer, c.Request)
	if err != nil {
		global.Logger.Error("客户端WS连接失败", zap.Error(err))
		return
	}

	go func() {
		time.Sleep(time.Second * 1)

		registry.ClientRegistry.Get(client.Code).Stop("此客户端已被其他客户端顶替，停止当前客户端，新连接的客户端IP：" + c.ClientIP())
		registry.ClientRegistry.Registry(client.Code, registry.NewV2Client(client.Code, handler))
		global.Cache.SetString(constant.CacheClientHeartbeatKey+client.Code, "online", time.Second*10)
		clientRegistry := registry.ClientRegistry.Get(client.Code)

		clientRegistry.Init()
		for _, host := range global.ClientHostFs.QueryAllFilter(func(host model.ClientHost) bool { return host.ClientCode == client.Code }) {
			clientRegistry.RunHost(host.Code, false)
		}
		for _, forward := range global.ClientForwardFs.QueryAllFilter(func(forward model.ClientForward) bool { return forward.ClientCode == client.Code }) {
			clientRegistry.RunForward(forward.Code, false)
		}
	}()
	go func() {
		upgrade.ReadLoop()
	}()
}
