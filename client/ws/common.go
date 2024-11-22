package ws

import (
	"sync"
)

var svrMap = make(map[string]int64)
var mux = &sync.RWMutex{}

func CheckUpdate(key string, tag int64) bool {
	mux.RLock()
	defer mux.RUnlock()
	value := svrMap[key]
	return tag > value
}

func Store(key string, value int64) {
	mux.Lock()
	defer mux.Unlock()
	svrMap[key] = value
}

func Remove(key string) {
	mux.Lock()
	defer mux.Unlock()
	delete(svrMap, key)
}
