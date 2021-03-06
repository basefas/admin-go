package middleware

import (
	"admin-go/internal/auth"
	"admin-go/internal/global"
	"admin-go/internal/utils"
	"admin-go/internal/utils/db"
	"strings"

	"github.com/gin-gonic/gin"
)

func Syslog() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		if token == "" {
			c.Next()
			return
		}
		userID, _ := auth.GetUID(token)
		path := c.Request.URL.Path
		method := c.Request.Method
		body := utils.GetRequestBody(c)
		clientIP := c.ClientIP()
		m := method == "POST" || method == "PUT" || method == "DELETE"
		p := !strings.Contains(path, "login")
		if m && p {
			log := global.OptLog{UserID: userID, Url: path, Method: method, Body: body, ClientIP: clientIP}
			db.Mysql.Create(&log)
		}
		c.Next()
	}
}
