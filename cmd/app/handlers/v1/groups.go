package v1

import (
	"go-admin/internal/groups"

	"github.com/gin-gonic/gin"
)

func GroupCreate(c *gin.Context) {
	var cg groups.CreateGroup
	if err := c.ShouldBindJSON(&cg); err != nil {
		Re(c, -1, InvalidArguments.Error(), nil)
		return
	}
	err := groups.Create(cg)
	if err != nil {
		Re(c, -1, err.Error(), nil)
	} else {
		Re(c, 0, "success", nil)
	}
}

func GroupGet(c *gin.Context) {
	groupID := c.Param("id")
	u, err := groups.Get(groupID)
	if err != nil {
		Re(c, -1, err.Error(), nil)
	} else {
		Re(c, 0, "success", u)
	}
}

func GroupUpdate(c *gin.Context) {
	groupID := c.Param("id")
	var ug groups.UpdateGroup
	if err := c.ShouldBindJSON(&ug); err != nil {
		Re(c, -1, err.Error(), nil)
		return
	}
	err := groups.Update(groupID, ug)
	if err != nil {
		Re(c, -1, err.Error(), nil)
	} else {
		Re(c, 0, "success", nil)
	}
}

func GroupDelete(c *gin.Context) {
	groupID := c.Param("id")
	err := groups.Delete(groupID)
	if err != nil {
		Re(c, -1, err.Error(), nil)
	} else {
		Re(c, 0, "success", nil)
	}
}

func GroupList(c *gin.Context) {
	gl, err := groups.List()
	if err != nil {
		Re(c, -1, err.Error(), nil)
	} else {
		Re(c, 0, "success", gl)
	}
}
