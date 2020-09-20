package v1

import (
	"go-admin/cmd/app/handlers"
	"go-admin/internal/users"

	"github.com/gin-gonic/gin"
)

func UsersCreate(c *gin.Context) {
	var cu users.CreateUser
	if err := c.ShouldBindJSON(&cu); err != nil {
		handlers.Re(c, -1, handlers.InvalidArguments.Error(), nil)
		return
	}

	err := users.Create(cu)
	if err != nil {
		handlers.Re(c, -1, err.Error(), nil)
	} else {
		handlers.Re(c, 0, "success", nil)
	}
}

func UsersGet(c *gin.Context) {
	userID := c.Param("id")

	u, err := users.Get(userID)
	if err != nil {
		handlers.Re(c, -1, err.Error(), nil)
	} else {
		handlers.Re(c, 0, "success", u)
	}
}

func UsersUpdate(c *gin.Context) {
	userID := c.Param("id")
	var uu users.UpdateUser
	if err := c.ShouldBindJSON(&uu); err != nil {
		handlers.Re(c, -1, handlers.InvalidArguments.Error(), nil)
		return
	}

	err := users.Update(userID, uu)
	if err != nil {
		handlers.Re(c, -1, err.Error(), nil)
	} else {
		handlers.Re(c, 0, "success", nil)
	}
}

func UsersDelete(c *gin.Context) {
	userID := c.Param("id")

	err := users.Delete(userID)
	if err != nil {
		handlers.Re(c, -1, err.Error(), nil)
	} else {
		handlers.Re(c, 0, "success", nil)
	}
}

func UsersList(c *gin.Context) {
	ul, err := users.List()
	if err != nil {
		handlers.Re(c, -1, err.Error(), nil)
	} else {
		handlers.Re(c, 0, "success", ul)
	}
}
