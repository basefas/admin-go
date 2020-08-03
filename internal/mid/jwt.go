package middleware

import (
	"go-admin/internal/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		code := 0
		token := c.GetHeader("token")
		if token == "" {
			code = -1
		} else {
			_, err := auth.ParseToken(token)
			if err != nil {
				code = -2
			}
		}

		if code != 0 {
			c.JSON(http.StatusOK, gin.H{
				"code":    code,
				"message": "Token Error",
				"data":    nil,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
