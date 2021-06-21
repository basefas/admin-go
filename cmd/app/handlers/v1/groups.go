package v1

import (
	"admin-go/cmd/app/handlers/http"
	"admin-go/internal/groups"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GroupCreate(c *gin.Context) {
	var cg groups.CreateGroup
	if err := c.ShouldBindJSON(&cg); err != nil {
		http.Re(c, -1, http.InvalidArguments.Error(), nil)
		return
	}
	err := groups.Create(cg)
	if err != nil {
		http.Re(c, -1, err.Error(), nil)
	} else {
		http.Re(c, 0, "success", nil)
	}
}

func GroupGet(c *gin.Context) {
	groupID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	g, err := groups.GetInfo(groupID)
	if err != nil {
		http.Re(c, -1, err.Error(), nil)
	} else {
		http.Re(c, 0, "success", g)
	}
}

func GroupUpdate(c *gin.Context) {
	groupID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var ug groups.UpdateGroup
	if err := c.ShouldBindJSON(&ug); err != nil {
		http.Re(c, -1, err.Error(), nil)
		return
	}
	err := groups.Update(groupID, ug)
	if err != nil {
		http.Re(c, -1, err.Error(), nil)
	} else {
		http.Re(c, 0, "success", nil)
	}
}

func GroupDelete(c *gin.Context) {
	groupID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	err := groups.Delete(groupID)
	if err != nil {
		http.Re(c, -1, err.Error(), nil)
	} else {
		http.Re(c, 0, "success", nil)
	}
}

func GroupList(c *gin.Context) {
	gl, err := groups.List()
	if err != nil {
		http.Re(c, -1, err.Error(), nil)
	} else {
		http.Re(c, 0, "success", gl)
	}
}
