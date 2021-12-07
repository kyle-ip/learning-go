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

	// 初始化需要绑定的服务提供者...

	// 将 HTTP 引擎初始化，并且作为服务提供者绑定到服务容器中。
	if engine, err := http.NewHttpEngine(); err == nil {
		container.Bind(&kernel.HadeKernelProvider{HttpEngine: engine})
	}

	// 运行 root 命令。
	console.RunCommand(container)

	//engine := gin.New()
	//
	//engine.Bind(&app.HadeAppProvider{})
	//engine.Bind(&demo.DemoProvider{})
	//
	//engine.Use(gin.Recovery())
	//engine.Use(middleware.Cost())
	//
	//hadeHttp.Routes(engine)
	//
	//server := &http.Server{
	//	Handler: engine,
	//	Addr:    ":8888",
	//}
	//// goroutine 启动服务。
	//go func() {
	//	server.ListenAndServe()
	//}()
	//
	//// 当前的goroutine等待信号量。
	//quit := make(chan os.Signal)
	//// 监控信号：SIGINT, SIGTERM, SIGQUIT。
	//signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	//// 阻塞当前 goroutine 等待信号。
	//<-quit
	//
	//// 调用 Server.Shutdown graceful 结束。
	//timeoutCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	//defer cancel()
	//
	//// if err := server.Shutdown(context.Background()); err != nil {
	//if err := server.Shutdown(timeoutCtx); err != nil {
	//	log.Fatal("Server Shutdown:", err)
	//}
}
