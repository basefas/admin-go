package v1

import (
	"admin-go/cmd/app/handlers/http"
	"admin-go/internal/roles"
	"strconv"

	"github.com/gin-gonic/gin"
)

func RoleCreate(c *gin.Context) {
	var cr roles.CreateRole
	if err := c.ShouldBindJSON(&cr); err != nil {
		http.Re(c, -1, err.Error(), nil)
		return
	}
	err := roles.Create(cr)
	if err != nil {
		http.Re(c, -1, err.Error(), nil)
	} else {
		http.Re(c, 0, "success", nil)
	}
}

func RoleGet(c *gin.Context) {
	roleID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	r, err := roles.GetInfo(roleID)
	if err != nil {
		http.Re(c, -1, err.Error(), nil)
	} else {
		http.Re(c, 0, "success", r)
	}
}

func RoleUpdate(c *gin.Context) {
	roleID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var ur roles.UpdateRole
	if err := c.ShouldBindJSON(&ur); err != nil {
		http.Re(c, -1, err.Error(), nil)
		return
	}
	err := roles.Update(roleID, ur)
	if err != nil {
		http.Re(c, -1, err.Error(), nil)
	} else {
		http.Re(c, 0, "success", nil)
	}
}

func RoleDelete(c *gin.Context) {
	roleID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	err := roles.Delete(roleID)
	if err != nil {
		http.Re(c, -1, err.Error(), nil)
	} else {
		http.Re(c, 0, "success", nil)
	}
}

func RoleList(c *gin.Context) {
	rl, err := roles.List()
	if err != nil {
		http.Re(c, -1, err.Error(), nil)
	} else {
		http.Re(c, 0, "success", rl)
	}
}

func RoleMenusList(c *gin.Context) {
	roleID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	rm, err := roles.GetRoleMenus(roleID)
	l := make([]uint64, 0)
	for _, menu := range rm {
		l = append(l, menu.MenuID)
	}
	if err != nil {
		http.Re(c, -1, err.Error(), nil)
	} else {
		http.Re(c, 0, "success", l)
	}
}

func RoleMenusUpdate(c *gin.Context) {
	roleID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var ml = make([]uint64, 0)
	if err := c.ShouldBindJSON(&ml); err != nil {
		http.Re(c, -1, err.Error(), nil)
		return
	}
	err := roles.UpdateRoleMenu(roleID, ml)

	if err != nil {
		http.Re(c, -1, err.Error(), nil)
	} else {
		http.Re(c, 0, "success", nil)
	}
}
