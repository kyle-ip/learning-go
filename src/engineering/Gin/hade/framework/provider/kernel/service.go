package kernel

import (
	"github.com/yipwinghong/hade/framework/gin"
	"net/http"
)

// HadeKernelService 引擎服务
type HadeKernelService struct {
	engine *gin.Engine
}

// NewHadeKernelService 初始化 web 引擎服务实例
func NewHadeKernelService(params ...interface{}) (interface{}, error) {
	httpEngine := params[0].(*gin.Engine)
	return &HadeKernelService{engine: httpEngine}, nil
}

// HttpEngine 返回 web 引擎
func (s *HadeKernelService) HttpEngine() http.Handler {
	return s.engine
}
