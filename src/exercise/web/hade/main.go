package main

import (
	"context"
	"github.com/yipwinghong/hade/framework/gin"
	"github.com/yipwinghong/hade/framework/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	core := gin.New()

	core.Use(gin.Recovery())
	core.Use(middleware.Cost())
	registerRouter(core)
	server := &http.Server{
		Handler: core,
		Addr:    ":8888",
	}
	// goroutine 启动服务。
	go func() {
		server.ListenAndServe()
	}()

	// 当前的goroutine等待信号量。
	quit := make(chan os.Signal)
	// 监控信号：SIGINT, SIGTERM, SIGQUIT。
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	// 阻塞当前 goroutine 等待信号。
	<-quit

	// 调用 Server.Shutdown graceful 结束。
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// if err := server.Shutdown(context.Background()); err != nil {
	if err := server.Shutdown(timeoutCtx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
}
