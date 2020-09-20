package v1

import (
	"go-admin/cmd/app/handlers"
	"go-admin/internal/auth"
	"go-admin/internal/menus"
	"go-admin/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func MenuCreate(c *gin.Context) {
	utils.GetRequestBody(c)
	var cm menus.CreateMenu
	if err := c.ShouldBindJSON(&cm); err != nil {
		//Re(c, -1, InvalidArguments.Error(), nil)
		handlers.Re(c, -1, err.Error(), nil)
		return
	}
	err := menus.Create(cm)
	if err != nil {
		handlers.Re(c, -1, err.Error(), nil)
	} else {
		handlers.Re(c, 0, "success", nil)
	}
}

func MenuGet(c *gin.Context) {
	menuID64, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	menuID := uint(menuID64)
	menuType := c.Query("type")
	var mi *menus.MenuInfo
	var err error
	if menuType == "tree" {
		mi, err = menus.GetTree(menuID)
	} else {
		mi, err = menus.Get(menuID)
	}
	if err != nil {
		handlers.Re(c, -1, err.Error(), nil)
	} else {
		handlers.Re(c, 0, "success", mi)
	}
}

func MenuUpdate(c *gin.Context) {
	menuID64, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	menuID := uint(menuID64)
	var um menus.UpdateMenu
	if err := c.ShouldBindJSON(&um); err != nil {
		handlers.Re(c, -1, err.Error(), nil)
		return
	}
	err := menus.Update(menuID, um)
	if err != nil {
		handlers.Re(c, -1, err.Error(), nil)
	} else {
		handlers.Re(c, 0, "success", nil)
	}
}

func MenuDelete(c *gin.Context) {
	menuID64, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	menuID := uint(menuID64)
	menuType := c.Query("type")
	var err error
	if menuType == "tree" {
		err = menus.DeleteTree(menuID)
	} else {
		err = menus.Delete(menuID)
	}

	if err != nil {
		handlers.Re(c, -1, err.Error(), nil)
	} else {
		handlers.Re(c, 0, "success", nil)
	}
}

func MenuList(c *gin.Context) {
	menuType := c.Query("type")
	var ml []menus.MenuInfo
	var err error

	switch menuType {
	case "tree":
		ml, err = menus.Tree()
	case "system":
		token := c.GetHeader("token")
		userID, _ := auth.GetUID(token)
		ml, err = menus.System(userID)
	default:
		ml, err = menus.List()
	}

	if err != nil {
		handlers.Re(c, -1, err.Error(), nil)
	} else {
		handlers.Re(c, 0, "success", ml)
	}
}
