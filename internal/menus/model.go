package menus

import (
	"admin-go/internal/global"
	"time"
)

type Menu struct {
	global.Model
	Name     string `json:"name" gorm:"NOT NULL"`
	Path     string `json:"path" gorm:"NOT NULL"`
	Method   string `json:"method" gorm:"NOT NULL"`
	Type     uint64 `json:"type" gorm:"type:uint;size:32;NOT NULL;"`
	Icon     string `json:"icon" gorm:"NOT NULL"`
	ParentID uint64 `json:"parent_id" gorm:"type:uint;size:32;NOT NULL;"`
	OrderID  uint64 `json:"order_id" gorm:"type:uint;size:32;NOT NULL;default:999999;"`
}

func (Menu) TableName() string {
	return "ag_menu"
}

type CreateMenu struct {
	Name     string `json:"name" binding:"required"`
	Path     string `json:"path" binding:"required"`
	Method   string `json:"method" binding:"-"`
	Type     uint64 `json:"type" binding:"required"`
	Icon     string `json:"icon" binding:"-"`
	ParentID uint64 `json:"parent_id" binding:"-"`
	OrderID  uint64 `json:"order_id" binding:"-"`
}

type UpdateMenu struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	Type     uint64 `json:"type"`
	Method   string `json:"method"`
	Icon     string `json:"icon"`
	ParentID uint64 `json:"parent_id"`
	OrderID  uint64 `json:"order_id"`
}

type MenuInfo struct {
	ID        uint64     `json:"id"`
	Name      string     `json:"name"`
	Path      string     `json:"path"`
	Type      uint64     `json:"type"`
	Method    string     `json:"method"`
	Icon      string     `json:"icon"`
	ParentID  uint64     `json:"parent_id"`
	OrderID   uint64     `json:"order_id"`
	CreatedAt time.Time  `json:"create_time"`
	UpdatedAt time.Time  `json:"update_time"`
	Children  []MenuInfo `json:"children" gorm:"-"`
	Funs      []MenuInfo `json:"funs" gorm:"-"`
}
