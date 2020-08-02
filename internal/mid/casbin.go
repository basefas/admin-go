package middleware

import (
	"go-admin/internal/auth"
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

func Casbin(e *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		allowed := auth.CheckPermission(c, e)
		if !allowed {
			c.JSON(http.StatusOK, gin.H{
				"code":    -1,
				"message": "auth error",
				"data":    nil,
			})
			c.Abort()
			return
		}
	}
}
