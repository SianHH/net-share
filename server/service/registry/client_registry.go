package registry

import (
	"sync"
)

type ClientInterface interface {
	RunHost(code string, force bool)
	DelHost(code string)
	RunForward(code string, force bool)
	DelForward(code string)
	RunTunnel(code string, force bool)
	DelTunnel(code string)
	Stop(reason string)
	Init()
	checkExec() bool
}

type clientRegistry struct {
	lock sync.RWMutex
	data map[string]ClientInterface
}

func (r *clientRegistry) Registry(code string, client ClientInterface) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.data[code] = client
}

func (r *clientRegistry) UnRegistry(code string) {
	r.lock.Lock()
	defer r.lock.Unlock()
	delete(r.data, code)
}

func (r *clientRegistry) Get(code string) ClientInterface {
	r.lock.RLock()
	defer r.lock.RUnlock()
	if value, ok := r.data[code]; ok {
		return value
	}
	return &V2Client{}
}

var ClientRegistry = clientRegistry{
	lock: sync.RWMutex{},
	data: make(map[string]ClientInterface),
}
