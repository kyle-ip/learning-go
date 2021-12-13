package main

import (
	"github.com/yipwinghong/hade/app/console"
	"github.com/yipwinghong/hade/app/http"
	"github.com/yipwinghong/hade/framework"
	"github.com/yipwinghong/hade/framework/provider/app"
	"github.com/yipwinghong/hade/framework/provider/kernel"
)

func main() {
	// 初始化服务容器。
	container := framework.NewHadeContainer()

	// 绑定 App 服务提供者。
	container.Bind(&app.HadeAppProvider{})

	// 将 HTTP 引擎初始化，并且作为服务提供者绑定到服务容器中。
	if engine, err := http.NewHttpEngine(); err == nil {
		container.Bind(&kernel.HadeKernelProvider{HttpEngine: engine})
	}

	// 运行 root 命令。
	console.RunCommand(container)
}
