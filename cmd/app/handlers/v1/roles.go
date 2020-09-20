package v1

import (
	"go-admin/cmd/app/handlers"
	"go-admin/internal/roles"

	"github.com/gin-gonic/gin"
)

func RoleCreate(c *gin.Context) {
	var cg roles.CreateRole
	if err := c.ShouldBindJSON(&cg); err != nil {
		handlers.Re(c, -1, handlers.InvalidArguments.Error(), nil)
		return
	}
	err := roles.Create(cg)
	if err != nil {
		handlers.Re(c, -1, err.Error(), nil)
	} else {
		handlers.Re(c, 0, "success", nil)
	}
}

func RoleGet(c *gin.Context) {
	roleID := c.Param("id")
	r, err := roles.Get(roleID)
	if err != nil {
		handlers.Re(c, -1, err.Error(), nil)
	} else {
		handlers.Re(c, 0, "success", r)
	}
}

func RoleUpdate(c *gin.Context) {
	roleID := c.Param("id")
	var ur roles.UpdateRole
	if err := c.ShouldBindJSON(&ur); err != nil {
		handlers.Re(c, -1, err.Error(), nil)
		return
	}
	err := roles.Update(roleID, ur)
	if err != nil {
		handlers.Re(c, -1, err.Error(), nil)
	} else {
		handlers.Re(c, 0, "success", nil)
	}
}

func RoleDelete(c *gin.Context) {
	roleID := c.Param("id")
	err := roles.Delete(roleID)
	if err != nil {
		handlers.Re(c, -1, err.Error(), nil)
	} else {
		handlers.Re(c, 0, "success", nil)
	}
}

func RoleList(c *gin.Context) {
	rl, err := roles.List()
	if err != nil {
		handlers.Re(c, -1, err.Error(), nil)
	} else {
		handlers.Re(c, 0, "success", rl)
	}
}

func RoleMenusList(c *gin.Context) {
	roleID := c.Param("id")
	rm, err := roles.GetRoleMenus(roleID)
	l := make([]uint, 0)
	for _, menu := range rm {
		l = append(l, menu.MenuID)
	}
	if err != nil {
		handlers.Re(c, -1, err.Error(), nil)
	} else {
		handlers.Re(c, 0, "success", l)
	}
}

func RoleMenusUpdate(c *gin.Context) {
	roleID := c.Param("id")
	var ml = make([]uint, 0)
	if err := c.ShouldBindJSON(&ml); err != nil {
		handlers.Re(c, -1, err.Error(), nil)
		return
	}
	err := roles.UpdateRoleMenu(roleID, ml)

	if err != nil {
		handlers.Re(c, -1, err.Error(), nil)
	} else {
		handlers.Re(c, 0, "success", nil)
	}
}
