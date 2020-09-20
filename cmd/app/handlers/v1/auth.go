package v1

import (
	"go-admin/cmd/app/handlers"
	"go-admin/internal/global"
	"go-admin/internal/users"
	"go-admin/internal/utils/db"

	"github.com/gin-gonic/gin"
)

func LogIn(c *gin.Context) {
	var u users.Login
	if err := c.ShouldBindJSON(&u); err != nil {
		handlers.Re(c, -1, err.Error(), nil)
		return
	}
	token, err := users.Token(u)
	if err != nil {
		log := global.AuthLog{Username: u.Username, ClientIP: c.ClientIP(), AuthStatus: 2}
		db.Mysql.Create(&log)
		handlers.Re(c, -1, err.Error(), nil)
	} else {
		log := global.AuthLog{Username: u.Username, ClientIP: c.ClientIP(), AuthStatus: 1}
		db.Mysql.Create(&log)
		handlers.Re(c, 0, "success", token)
	}
}
