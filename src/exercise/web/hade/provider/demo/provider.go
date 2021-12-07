package demo

import (
	"fmt"
	"github.com/yipwinghong/hade/framework"
)

// DemoServiceProvider 服务提供方。
type DemoServiceProvider struct {
}

// Name 方法直接将服务对应的字符串凭证返回，hade.demo。
func (sp *DemoServiceProvider) Name() string {
	return Key
}

// Register 注册初始化服务实例。
func (sp *DemoServiceProvider) Register(c framework.Container) framework.NewInstance {
	return NewDemoService
}

// IsDefer 是否延迟实例化，此处实例化延迟到第一次 make 时。
func (sp *DemoServiceProvider) IsDefer() bool {
	return true
}

// Params 实例化的参数，此处表示在 NewDemoService 中只有参数 container。
func (sp *DemoServiceProvider) Params(c framework.Container) []interface{} {
	return []interface{}{c}
}

// Boot 只打印日志。
func (sp *DemoServiceProvider) Boot(c framework.Container) error {
	fmt.Println("demo service boot")
	return nil
}
