package framework

import (
	"errors"
	"fmt"
	"sync"
)

// Container 服务容器，提供绑定服务和获取服务的功能。
type Container interface {

	// Bind 绑定一个服务提供者，如果关键字凭证已存在会进行替换操作，返回 error。
	Bind(provider ServiceProvider) error

	// IsBind 关键字凭证是否已经绑定服务提供者。
	IsBind(key string) bool

	// Make 根据关键字凭证获取服务。
	Make(key string) (interface{}, error)

	// MustMake 根据关键字凭证获取服务，如果关键字凭证未绑定服务提供者会 panic。
	// 在使用这个接口时请保证服务容器已为关键字凭证绑定服务提供者。
	MustMake(key string) interface{}

	// MakeNew 根据关键字凭证获取服务，该服务非单例。
	// 根据服务提供者注册的启动函数和传递的 params 参数实例化（在需要为不同参数启动不同实例时适用）。
	MakeNew(key string, params []interface{}) (interface{}, error)
}

// HadeContainer 服务容器的具体实现。
type HadeContainer struct {
	Container
	// providers 注册服务提供者，key 为字符串凭证。
	providers map[string]ServiceProvider

	// instance 具体的服务实例，key 为字符串凭证。
	instances map[string]interface{}

	// lock 锁定容器变更操作。
	lock sync.RWMutex
}

// NewHadeContainer 创建服务容器
func NewHadeContainer() *HadeContainer {
	return &HadeContainer{
		providers: map[string]ServiceProvider{},
		instances: map[string]interface{}{},
		lock:      sync.RWMutex{},
	}
}

// PrintProviders 输出服务容器中注册的关键字
func (hade *HadeContainer) PrintProviders() []string {
	ret := make([]string, len(hade.providers))
	// ret := []string{}
	for _, provider := range hade.providers {
		name := provider.Name()

		line := fmt.Sprint(name)
		ret = append(ret, line)
	}
	return ret
}

// Bind 绑定服务容器和关键字。
func (hade *HadeContainer) Bind(provider ServiceProvider) error {
	hade.lock.Lock()
	defer hade.lock.Unlock()

	key := provider.Name()
	hade.providers[key] = provider

	if provider.IsDefer() == false {
		if err := provider.Boot(hade); err != nil {
			return err
		}
		// 实例化方法。
		params := provider.Params(hade)
		method := provider.Register(hade)
		instance, err := method(params...)
		if err != nil {
			return errors.New(err.Error())
		}
		hade.instances[key] = instance
	}
	return nil
}

func (hade *HadeContainer) IsBind(key string) bool {
	return hade.findServiceProvider(key) != nil
}

func (hade *HadeContainer) findServiceProvider(key string) ServiceProvider {
	hade.lock.RLock()
	defer hade.lock.RUnlock()
	if sp, ok := hade.providers[key]; ok {
		return sp
	}
	return nil
}

// Make 实例化（单例）。
func (hade *HadeContainer) Make(key string) (interface{}, error) {
	return hade.make(key, nil, false)
}

// MustMake 强制实例化。
func (hade *HadeContainer) MustMake(key string) interface{} {
	serv, err := hade.make(key, nil, false)
	if err != nil {
		panic(err)
	}
	return serv
}

// MakeNew 实例化。
func (hade *HadeContainer) MakeNew(key string, params []interface{}) (interface{}, error) {
	return hade.make(key, params, true)
}

// make 执行服务实例化。
func (hade *HadeContainer) make(key string, params []interface{}, forceNew bool) (interface{}, error) {
	hade.lock.RLock()
	defer hade.lock.RUnlock()

	// 查询是否已经注册了这个服务提供者，如果没有注册则返回错误。
	sp := hade.findServiceProvider(key)
	if sp == nil {
		return nil, errors.New("contract " + key + " have not register")
	}

	// 强制重新实例化。
	if forceNew {
		return hade.newInstance(sp, params)
	}

	// 如果容器中已经实例化，就直接使用容器中的实例，否则执行实例化、添加到实例集合中。
	if ins, ok := hade.instances[key]; ok {
		return ins, nil
	}
	inst, err := hade.newInstance(sp, nil)
	if err != nil {
		return nil, err
	}
	hade.instances[key] = inst
	return inst, nil
}

// newInstance
func (hade *HadeContainer) newInstance(sp ServiceProvider, params []interface{}) (interface{}, error) {
	if err := sp.Boot(hade); err != nil {
		return nil, err
	}
	if params == nil {
		params = sp.Params(hade)
	}
	method := sp.Register(hade)
	ins, err := method(params...)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return ins, err
}
