package v1

import (
	"fmt"
	"go-admin/internal/auth"
	"go-admin/internal/users"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func UserGet(c *gin.Context) {
	id := c.Param("id")
	token := c.GetHeader("token")
	userID, _ := auth.GetUID(token)

	if !checkID(id, userID) {
		Re(c, -1, errors.New("auth error.").Error(), nil)
		return
	}

	u, err := users.Get(id)

	if err != nil {
		Re(c, -1, err.Error(), nil)
	} else {
		Re(c, 0, "success", u)
	}
}

func UserUpdate(c *gin.Context) {
	id := c.Param("id")
	token := c.GetHeader("token")
	userID, _ := auth.GetUID(token)

	if !checkID(id, userID) {
		Re(c, -1, AuthError.Error(), nil)
		return
	}

	var uu users.UpdateUser
	if err := c.ShouldBindJSON(&uu); err != nil {
		Re(c, -1, InvalidArguments.Error(), nil)
		return
	}

	err := users.Update(id, uu)

	if err != nil {
		Re(c, -1, err.Error(), nil)
	} else {
		Re(c, 0, "success", nil)
	}
}

func checkID(id string, userID uint) bool {
	return id == fmt.Sprint(userID)
}
