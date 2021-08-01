package router

import (
    "github.com/gin-gonic/gin"
    "week02/internal/service"
)

func New() *gin.Engine {
    r := gin.Default()
    r.GET("/user/get", service.GetUser())
    return r
}
