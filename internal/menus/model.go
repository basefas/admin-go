package menus

import (
	"go-admin/internal/global"
	"time"
)

type Menu struct {
	global.Model
	MenuName string `json:"name" gorm:"NOT NULL"`
	MenuPath string `json:"path" gorm:"NOT NULL"`
	Method   string `json:"method" gorm:"NOT NULL"`
	MenuType uint   `json:"menu_type" gorm:"NOT NULL"`
	Icon     string `json:"icon" gorm:"NOT NULL"`
	ParentID uint   `json:"parent_id" gorm:"NOT NULL"`
	OrderID  uint   `json:"order_id" gorm:"NOT NULL" gorm:"default:999999"`
}

type CreateMenu struct {
	MenuName string `json:"name" binding:"required"`
	MenuPath string `json:"path" binding:"required"`
	Method   string `json:"method" binding:"-"`
	MenuType uint   `json:"menu_type" binding:"required"`
	Icon     string `json:"icon" binding:"-"`
	ParentID uint   `json:"parent_id" binding:"-"`
	OrderID  uint   `json:"order_id" binding:"-"`
}

type UpdateMenu struct {
	MenuName string `json:"name"`
	MenuPath string `json:"path"`
	MenuType uint   `json:"menu_type"`
	Method   string `json:"method"`
	Icon     string `json:"icon"`
	ParentID uint   `json:"parent_id"`
	OrderID  uint   `json:"order_id"`
}

type MenuInfo struct {
	ID        uint        `json:"id"`
	MenuName  string      `json:"name"`
	MenuPath  string      `json:"path"`
	MenuType  uint        `json:"menu_type"`
	Method    string      `json:"method"`
	Icon      string      `json:"icon"`
	ParentID  uint        `json:"parent_id"`
	OrderID   uint        `json:"order_id"`
	CreatedAt time.Time   `json:"create_time"`
	UpdatedAt time.Time   `json:"update_time"`
	Children  []*MenuInfo `json:"children"`
	Funs      []*MenuInfo `json:"funs"`
}

type RoleMenu struct {
	global.Model
	RoleID uint `json:"role_id" gorm:"NOT NULL"`
	MenuID uint `json:"menu_id" gorm:"NOT NULL"`
}
