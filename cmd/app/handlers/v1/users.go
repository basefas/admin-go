package v1

import (
	"admin-go/cmd/app/handlers/http"
	"admin-go/internal/users"
	"strconv"

	"github.com/gin-gonic/gin"
)

func UsersCreate(c *gin.Context) {
	var cu users.CreateUser
	if err := c.ShouldBindJSON(&cu); err != nil {
		http.Re(c, -1, http.InvalidArguments.Error(), nil)
		return
	}
	err := users.Create(cu)
	if err != nil {
		http.Re(c, -1, err.Error(), nil)
	} else {
		http.Re(c, 0, "success", nil)
	}
}

func UsersGet(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	u, err := users.GetInfo(userID)
	if err != nil {
		http.Re(c, -1, err.Error(), nil)
	} else {
		http.Re(c, 0, "success", u)
	}
}

func UsersUpdate(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var uu users.UpdateUser
	if err := c.ShouldBindJSON(&uu); err != nil {
		http.Re(c, -1, http.InvalidArguments.Error(), nil)
		return
	}

	err := users.Update(userID, uu)
	if err != nil {
		http.Re(c, -1, err.Error(), nil)
	} else {
		http.Re(c, 0, "success", nil)
	}
}

func UsersDelete(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	err := users.Delete(userID)
	if err != nil {
		http.Re(c, -1, err.Error(), nil)
	} else {
		http.Re(c, 0, "success", nil)
	}
}

func UsersList(c *gin.Context) {
	ul, err := users.List()
	if err != nil {
		http.Re(c, -1, err.Error(), nil)
	} else {
		http.Re(c, 0, "success", ul)
	}
}
