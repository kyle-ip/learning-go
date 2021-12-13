package middleware

import (
	"api-service/pkg/errno"
	"api-service/pkg/token"
	"api-service/util"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse the json web token.
		if _, err := token.ParseRequest(c); err != nil {
			util.SendResponse(c, errno.ErrTokenInvalid, nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
