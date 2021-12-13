package gin

import "github.com/yipwinghong/hade/framework"

func (engine *Engine) SetContainer(container framework.Container) {
	engine.container = container
}

// Bind engine 实现 container 的绑定封装
func (engine *Engine) Bind(provider framework.ServiceProvider) error {
	return engine.container.Bind(provider)
}

// IsBind 关键字凭证是否已经绑定服务提供者
func (engine *Engine) IsBind(key string) bool {
	return engine.container.IsBind(key)
}
