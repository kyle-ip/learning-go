package http

import (
	"github.com/yipwinghong/hade/framework/gin"
)

// NewHttpEngine 创建了一个绑定了路由的Web引擎
func NewHttpEngine() (*gin.Engine, error) {

	// 设置为 Release，为的是默认在启动中不输出调试信息
	gin.SetMode(gin.ReleaseMode)

	// 默认启动一个Web引擎
	r := gin.Default()

	// 业务绑定路由操作
	Routes(r)
	// 返回绑定路由后的Web引擎
	return r, nil
}
