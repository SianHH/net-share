package component

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"net-share/pkg/signal"
	"net-share/server/framework/hook"
	"net-share/server/global"
	"net/http"
	"os"
	"time"
)

func InitServer() {
	hook.Run()

	server := &http.Server{
		Addr:              global.App.Addr,
		Handler:           engine,
		MaxHeaderBytes:    0,
		WriteTimeout:      time.Second * 30,
		ReadHeaderTimeout: time.Second * 30,
		ReadTimeout:       time.Second * 30,
	}

	// 资源释放
	var free = func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		_ = server.Shutdown(ctx)
		global.ClientFs.Close()
		global.ClientHostFs.Close()
		global.ClientHostDomainFs.Close()
		global.ClientForwardFs.Close()
		_ = global.Cache.Sync()
		global.Logger.Info("释放资源，运行结束")
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			global.Logger.Error("启动失败", zap.Error(err))
			free()
			os.Exit(1)
		}
	}()

	time.Sleep(time.Second * 1)
	global.Logger.Info("服务启动成功 " + global.App.Addr)
	signal.Func(free)
}
