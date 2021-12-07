package main

import (
	"github.com/yipwinghong/hade/framework/gin"
	"time"
)

func UserLoginController(c *gin.Context) {
	foo, _ := c.DefaultQueryString("foo", "def")
	// 等待 10s 才继续执行。
	time.Sleep(10 * time.Second)
	c.ISetOkStatus().IJson("ok, UserLoginController: " + foo)
}
