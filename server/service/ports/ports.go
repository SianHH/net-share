package ports

import (
	"errors"
	"github.com/patrickmn/go-cache"
	"net-share/pkg/utils"
	"net-share/server/constant"
	"net-share/server/framework/hook"
	"net-share/server/global"
	"net-share/server/model"
	"strconv"
	"strings"
	"sync"
	"time"
)

func init() {
	hook.AddServerBeforeHookFunc(func() {
		go func() {
			for {
				global.Logger.Info("开始整理可用端口")
				PortFunc()
				global.Logger.Info("结束整理可用端口")
				time.Sleep(time.Minute * 60)
			}
		}()
	})
}

var portLock = &sync.Mutex{}

func PortFunc() {
	portLock.Lock()
	defer portLock.Unlock()

	var portList = utils.Map(global.ClientForwardFs.QueryAll(), func(t model.ClientForward) (string, bool) {
		return t.Port, true
	})
	allowPort := analysisPorts(portList)
	global.Cache.SetStruct(constant.CacheNodeAllowPortKey, allowPort, cache.NoExpiration)
}

func PullPort() (string, error) {
	portLock.Lock()
	defer portLock.Unlock()
	var portList []string
	if err := global.Cache.GetStruct(constant.CacheNodeAllowPortKey, &portList); err != nil {
		return "", err
	}
	if len(portList) == 0 {
		return "", errors.New("端口资源不足")
	}
	result := portList[0]
	portList = portList[1:]
	global.Cache.SetStruct(constant.CacheNodeAllowPortKey, portList, cache.NoExpiration)
	return result, nil
}
func PushPort(port string) {
	portLock.Lock()
	defer portLock.Unlock()
	var portList []string
	if err := global.Cache.GetStruct(constant.CacheNodeAllowPortKey, &portList); err != nil {
		return
	}
	portList = append(portList, port)
	global.Cache.SetStruct(constant.CacheNodeAllowPortKey, portList, cache.NoExpiration)
}

// 解析端口
func analysisPorts(excludePort []string) (result []string) {
	list := []string{}
	var excludePortMap = make(map[string]bool)
	for _, port := range excludePort {
		excludePortMap[port] = true
	}
	for _, v1 := range global.App.Ports {
		if v1 == "" {
			continue
		}
		if _, err := strconv.Atoi(v1); err == nil {
			list = append(list, v1)
		}
		portGroup := strings.Split(v1, "-")
		if len(portGroup) != 2 {
			continue
		}
		start, err := strconv.Atoi(portGroup[0])
		if err != nil {
			continue
		}
		end, err := strconv.Atoi(portGroup[1])
		if err != nil {
			continue
		}
		if start >= end {
			continue
		}
		for {
			if start > end {
				break
			}
			list = append(list, strconv.Itoa(start))
			start++
		}
	}
	for _, item := range list {
		if excludePortMap[item] {
			continue
		}
		result = append(result, item)
	}
	return result
}
