package main

import (
	"github.com/yipwinghong/hade/framework/gin"
	"github.com/yipwinghong/hade/framework/middleware"
)

// registerRouter 注册路由规则
func registerRouter(core *gin.Engine) {
	// 静态路由 + HTTP 方法匹配。
	core.GET("/user/login", middleware.Test3(), UserLoginController)

	// 批量通用前缀。
	subjectApi := core.Group("/subject")

	// 动态路由。
	{
		subjectApi.Use(middleware.Test3())
		subjectApi.DELETE("/:id", SubjectDelController)
		subjectApi.PUT("/:id", SubjectUpdateController)
		subjectApi.GET("/:id", middleware.Test3(), SubjectGetController)
		subjectApi.GET("/list/all", SubjectListController)

		subjectInnerApi := subjectApi.Group("/info")
		{
			subjectInnerApi.GET("/name", SubjectNameController)
		}
	}
}
