package v1

import (
	"admin-go/cmd/app/handlers/http"
	"admin-go/internal/global"
	"admin-go/internal/users"
	"admin-go/internal/utils/db"

	"github.com/gin-gonic/gin"
)

func LogIn(c *gin.Context) {
	var u users.Login
	if err := c.ShouldBindJSON(&u); err != nil {
		http.Re(c, -1, err.Error(), nil)
		return
	}
	token, err := users.Token(u)
	if err != nil {
		log := global.AuthLog{Username: u.Username, ClientIP: c.ClientIP(), AuthStatus: 2}
		db.Mysql.Create(&log)
		http.Re(c, -1, err.Error(), nil)
	} else {
		log := global.AuthLog{Username: u.Username, ClientIP: c.ClientIP(), AuthStatus: 1}
		db.Mysql.Create(&log)
		http.Re(c, 0, "success", token)
	}
}
