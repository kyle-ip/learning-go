package main

import (
	"coredemo/framework"
	"coredemo/framework/middleware"
)

// registerRouter 注册路由规则
func registerRouter(core *framework.Core) {
	// 静态路由 + HTTP 方法匹配。
	core.Get("/user/login", middleware.Test3(), UserLoginController)

	// 批量通用前缀。
	subjectApi := core.Group("/subject")

	// 动态路由。
	{
		subjectApi.Use(middleware.Test3())
		subjectApi.Delete("/:id", SubjectDelController)
		subjectApi.Put("/:id", SubjectUpdateController)
		subjectApi.Get("/:id", middleware.Test3(), SubjectGetController)
		subjectApi.Get("/list/all", SubjectListController)

		subjectInnerApi := subjectApi.Group("/info")
		{
			subjectInnerApi.Get("/name", SubjectNameController)
		}
	}
}
