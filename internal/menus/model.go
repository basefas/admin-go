package menus

import (
	"admin-go/internal/global"
	"time"
)

type Menu struct {
	global.Model
	MenuName string `json:"name" gorm:"NOT NULL"`
	MenuPath string `json:"path" gorm:"NOT NULL"`
	Method   string `json:"method" gorm:"NOT NULL"`
	MenuType uint64 `json:"menu_type" gorm:"type:uint;size:32;NOT NULL;"`
	Icon     string `json:"icon" gorm:"NOT NULL"`
	ParentID uint64 `json:"parent_id" gorm:"type:uint;size:32;NOT NULL;"`
	OrderID  uint64 `json:"order_id" gorm:"type:uint;size:32;NOT NULL;default:999999;"`
}

type CreateMenu struct {
	MenuName string `json:"name" binding:"required"`
	MenuPath string `json:"path" binding:"required"`
	Method   string `json:"method" binding:"-"`
	MenuType uint64 `json:"menu_type" binding:"required"`
	Icon     string `json:"icon" binding:"-"`
	ParentID uint64 `json:"parent_id" binding:"-"`
	OrderID  uint64 `json:"order_id" binding:"-"`
}

type UpdateMenu struct {
	MenuName string `json:"name"`
	MenuPath string `json:"path"`
	MenuType uint64 `json:"menu_type"`
	Method   string `json:"method"`
	Icon     string `json:"icon"`
	ParentID uint64 `json:"parent_id"`
	OrderID  uint64 `json:"order_id"`
}

type MenuInfo struct {
	ID        uint64     `json:"id"`
	MenuName  string     `json:"name"`
	MenuPath  string     `json:"path"`
	MenuType  uint64     `json:"menu_type"`
	Method    string     `json:"method"`
	Icon      string     `json:"icon"`
	ParentID  uint64     `json:"parent_id"`
	OrderID   uint64     `json:"order_id"`
	CreatedAt time.Time  `json:"create_time"`
	UpdatedAt time.Time  `json:"update_time"`
	Children  []MenuInfo `json:"children" gorm:"-"`
	Funs      []MenuInfo `json:"funs" gorm:"-"`
}

type RoleMenu struct {
	global.Model
	RoleID uint64 `json:"role_id" gorm:"type:uint;size:32;NOT NULL;"`
	MenuID uint64 `json:"menu_id" gorm:"type:uint;size:32;NOT NULL;"`
}
