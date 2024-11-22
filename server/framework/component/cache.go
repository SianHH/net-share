package component

import (
	"net-share/pkg/cache"
	"net-share/server/global"
	"time"
)

func InitCache() {
	global.Cache = cache.NewCache("data/cache.db")
	go func() {
		for {
			time.Sleep(time.Minute * 10)
			_ = global.Cache.Sync()
		}
	}()

	//global.Cache.OnEvicted(func(key string, value interface{}) {
	//
	//})
}
