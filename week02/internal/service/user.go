package service

import (
    "errors"
    "github.com/gin-gonic/gin"
    "net/http"
    "strconv"
    "week02/internal/dao"
)

func GetUser() gin.HandlerFunc {
    return func(c *gin.Context) {
        query := c.Query("id")
        aId, err := strconv.ParseInt(query, 10, 64)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"message": "参数校验失败"})
            return
        }
        user, err := dao.GetUser(aId)
        if errors.Is(err, dao.ErrNoRows) {
            c.JSON(http.StatusBadRequest, gin.H{"message": "找不到用户"})
            return
        }
        c.JSON(http.StatusOK, gin.H{
            "id":       user.Id,
            "username": user.Username,
        })
    }
}
