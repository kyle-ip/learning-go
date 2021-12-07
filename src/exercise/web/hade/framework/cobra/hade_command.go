package cobra

import "github.com/yipwinghong/hade/framework"

// SetContainer 设置容器
func (c *Command) SetContainer(container framework.Container) {
	c.container = container
}

// GetContainer 获取容器
func (c *Command) GetContainer() framework.Container {
	return c.Root().container
}
