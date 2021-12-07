package demo

import (
	"fmt"
	"github.com/yipwinghong/hade/framework"
)

// DemoService 具体的接口实例
type DemoService struct {
	Service

	// 参数
	c framework.Container
}

// NewDemoService 初始化实例
func NewDemoService(params ...interface{}) (interface{}, error) {
	// 这里需要将参数展开
	c := params[0].(framework.Container)

	fmt.Println("new demo service")
	// 返回实例
	return &DemoService{c: c}, nil
}

// GetFoo 实现接口
func (s *DemoService) GetFoo() Foo {
	return Foo{
		Name: "i am foo",
	}
}
