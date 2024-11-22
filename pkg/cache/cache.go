package cache

import (
	"encoding/json"
	"errors"
	"github.com/patrickmn/go-cache"
	"time"
)

type Cache struct {
	*cache.Cache
	file string
}

func (c *Cache) Sync() error {
	return c.Cache.SaveFile(c.file)
}

func NewCache(file string) *Cache {
	c := cache.New(5*time.Minute, 10*time.Minute)
	_ = c.LoadFile(file)
	go func() {
		for {
			// 定期实例化
			time.Sleep(time.Minute * 10)
			_ = c.SaveFile("cache.db")
		}
	}()
	return &Cache{c, file}
}

func (c *Cache) SetStruct(key string, data any, duration time.Duration) {
	marshal, _ := json.Marshal(data)
	c.Set(key, marshal, duration)
}

func (c *Cache) GetStruct(key string, data any) error {
	val, ok := c.Get(key)
	if !ok {
		return errors.New("not key")
	}
	return json.Unmarshal(val.([]byte), data)
}

func (c *Cache) GetString(key string) (value string) {
	val, ok := c.Get(key)
	if !ok {
		return ""
	}
	str, ok := val.(string)
	if ok {
		return str
	}
	return ""
}

func (c *Cache) SetString(key string, value string, duration time.Duration) {
	c.Set(key, value, duration)
}
