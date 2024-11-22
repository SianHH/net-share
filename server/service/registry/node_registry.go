package registry

import (
	"github.com/go-gost/core/logger"
	"github.com/go-gost/x/config"
	"github.com/go-gost/x/config/parsing"
	"github.com/go-gost/x/config/parsing/auth"
	"github.com/go-gost/x/config/parsing/ingress"
	xservice "github.com/go-gost/x/config/parsing/service"
	xlogger "github.com/go-gost/x/logger"
	"github.com/go-gost/x/registry"
	"go.uber.org/zap"
	"log"
	"net-share/server/framework/hook"
	"net-share/server/global"
	"os"
	"time"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
	logger.SetDefault(xlogger.NewLogger(xlogger.LevelOption(logger.FatalLevel)))
	parsing.BuildDefaultTLSConfig(nil)

	hook.AddServerBeforeHookFunc(func() {
		runHostServer()
		runForwardServer()
		UpdateIngress()
		UpdateAuthers()
	})
}

func UpdateAuthers() {
	var authers []*config.AuthConfig
	for _, host := range global.ClientHostFs.QueryAll() {
		authers = append(authers, &config.AuthConfig{
			Username: host.AuthUser,
			Password: host.AuthPwd,
		})
	}
	for _, forward := range global.ClientForwardFs.QueryAll() {
		authers = append(authers, &config.AuthConfig{
			Username: forward.AuthUser,
			Password: forward.AuthPwd,
		})
	}

	for _, tunnel := range global.ClientTunnelFs.QueryAll() {
		authers = append(authers, &config.AuthConfig{
			Username: tunnel.AuthUser,
			Password: tunnel.AuthPwd,
		})
	}

	auther := auth.ParseAuther(&config.AutherConfig{
		Name:  "authers",
		Auths: authers,
	})
	registry.AutherRegistry().Unregister("authers")
	_ = registry.AutherRegistry().Register("authers", auther)
}

func UpdateIngress() {
	var rules []*config.IngressRuleConfig
	for _, host := range global.ClientHostFs.QueryAll() {
		rules = append(rules, &config.IngressRuleConfig{
			Hostname: host.DomainPrefix + "." + global.App.Domain,
			Endpoint: host.Code,
		})
	}

	for _, tunnel := range global.ClientTunnelFs.QueryAll() {
		rules = append(rules, &config.IngressRuleConfig{
			Hostname: tunnel.DomainPrefix + "." + global.App.Domain,
			Endpoint: "$" + tunnel.Code,
		})
	}

	parseIngress := ingress.ParseIngress(&config.IngressConfig{
		Name:  "ingress",
		Rules: rules,
	})
	registry.IngressRegistry().Unregister("ingress")
	_ = registry.IngressRegistry().Register("ingress", parseIngress)
}

func runHostServer() {
	if global.App.HostPort == "" || global.App.Domain == "" || global.App.Entrypoint == "" {
		global.Logger.Warn("未配置host-port、domain、entrypoint，不启用域名解析")
		return
	}
	parseService, err := xservice.ParseService(&config.ServiceConfig{
		Name: "host",
		Addr: ":" + global.App.HostPort,
		Handler: &config.HandlerConfig{
			Type: "tunnel",
			Metadata: map[string]any{
				"entrypoint": ":" + global.App.Entrypoint,
				"ingress":    "ingress",
				"sniffing":   true,
			},
			Auther: "authers",
		},
		Listener: &config.ListenerConfig{
			Type: "tls",
		},
	})
	if err != nil {
		global.Logger.Warn("加载HostServer配置错误", zap.Error(err))
		return
	}
	_ = registry.ServiceRegistry().Register("hostServer", parseService)
	go func() {
		if err := parseService.Serve(); err != nil {
			global.Logger.Error("启动HostServer失败", zap.Error(err))
			os.Exit(1)
		}
	}()
	time.Sleep(time.Second * 1)
	global.Logger.Info("启动HostServer成功")
}

func runForwardServer() {
	if global.App.ForwardPort == "" || len(global.App.Ports) == 0 {
		global.Logger.Warn("未配置forward-port或者ports，不启用端口转发")
		return
	}
	parseService, err := xservice.ParseService(&config.ServiceConfig{
		Name: "forward",
		Addr: ":" + global.App.ForwardPort,
		Handler: &config.HandlerConfig{
			Type: "relay",
			Metadata: map[string]any{
				"nodelay": true,
				"bind":    true,
			},
			Auther: "authers",
		},
		Listener: &config.ListenerConfig{
			Type: "tls",
		},
	})
	if err != nil {
		global.Logger.Warn("加载ForwardServer配置错误", zap.Error(err))
		return
	}
	_ = registry.ServiceRegistry().Register("forwardServer", parseService)
	go func() {
		if err := parseService.Serve(); err != nil {
			global.Logger.Error("启动ForwardServer失败", zap.Error(err))
			os.Exit(1)
		}
	}()
	time.Sleep(time.Second * 1)
	global.Logger.Info("启动ForwardServer成功")
}
