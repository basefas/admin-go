package v1

import (
	"admin-go/cmd/app/handlers/http"
	"admin-go/internal/auth"
	"admin-go/internal/menus"
	"strconv"

	"github.com/gin-gonic/gin"
)

func MenuCreate(c *gin.Context) {
	var cm menus.CreateMenu
	if err := c.ShouldBindJSON(&cm); err != nil {
		http.Re(c, -1, err.Error(), nil)
		return
	}
	err := menus.Create(cm)
	if err != nil {
		http.Re(c, -1, err.Error(), nil)
	} else {
		http.Re(c, 0, "success", nil)
	}
}

func MenuGet(c *gin.Context) {
	menuID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	menuType := c.Query("type")
	var mi menus.MenuInfo
	var err error
	if menuType == "tree" {
		mi, err = menus.GetTree(menuID)
	} else {
		mi, err = menus.GetInfo(menuID)
	}
	if err != nil {
		http.Re(c, -1, err.Error(), nil)
	} else {
		http.Re(c, 0, "success", mi)
	}
}

func MenuUpdate(c *gin.Context) {
	menuID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var um menus.UpdateMenu
	if err := c.ShouldBindJSON(&um); err != nil {
		http.Re(c, -1, err.Error(), nil)
		return
	}
	err := menus.Update(menuID, um)
	if err != nil {
		http.Re(c, -1, err.Error(), nil)
	} else {
		http.Re(c, 0, "success", nil)
	}
}

func MenuDelete(c *gin.Context) {
	menuID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	err := menus.Delete(menuID)
	if err != nil {
		http.Re(c, -1, err.Error(), nil)
	} else {
		http.Re(c, 0, "success", nil)
	}
}

func MenuList(c *gin.Context) {
	menuType := c.Query("type")
	ml := make([]menus.MenuInfo, 0)
	var err error

	switch menuType {
	case "tree":
		ml, err = menus.Tree()
	case "system":
		userID, _ := auth.GetUID(c.GetHeader("token"))
		ml, err = menus.System(userID)
	default:
		ml, err = menus.List()
	}

	if err != nil {
		http.Re(c, -1, err.Error(), nil)
	} else {
		http.Re(c, 0, "success", ml)
	}
}
