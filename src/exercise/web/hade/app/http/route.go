package http

import (
	"github.com/yipwinghong/hade/app/http/module/demo"
	"github.com/yipwinghong/hade/framework/gin"
)

func Routes(r *gin.Engine) {

	r.Static("/dist/", "./dist/")

	demo.Register(r)
}
