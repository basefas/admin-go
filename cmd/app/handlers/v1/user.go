package v1

import (
	"admin-go/cmd/app/handlers/http"
	"admin-go/internal/auth"
	"admin-go/internal/users"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func UserGet(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	token := c.GetHeader("token")
	userID, _ := auth.GetUID(token)

	if !(id == userID) {
		http.Re(c, -1, errors.New("auth error.").Error(), nil)
		return
	}

	u, err := users.Get(id)

	if err != nil {
		http.Re(c, -1, err.Error(), nil)
	} else {
		http.Re(c, 0, "success", u)
	}
}

func UserUpdate(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	token := c.GetHeader("token")
	userID, _ := auth.GetUID(token)

	if !(id == userID) {
		http.Re(c, -1, http.AuthError.Error(), nil)
		return
	}

	var uu users.UpdateUser
	if err := c.ShouldBindJSON(&uu); err != nil {
		http.Re(c, -1, http.InvalidArguments.Error(), nil)
		return
	}
	err := users.Update(id, uu)

	if err != nil {
		http.Re(c, -1, err.Error(), nil)
	} else {
		http.Re(c, 0, "success", nil)
	}
}
