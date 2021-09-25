package roles

import (
	"admin-go/internal/global"
	"time"
)

type Role struct {
	global.Model
	Name string `json:"name" gorm:"NOT NULL"`
}

func (Role) TableName() string {
	return "ag_role"
}

type CreateRole struct {
	Name string `json:"name" binding:"required"`
}

type UpdateRole struct {
	Name string `json:"name"`
}

type RoleInfo struct {
	ID        uint64    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"create_time"`
	UpdatedAt time.Time `json:"update_time"`
}

type RoleMenu struct {
	global.Model
	RoleID uint64 `json:"role_id" gorm:"type:uint;size:32;NOT NULL;"`
	MenuID uint64 `json:"menu_id" gorm:"type:uint;size:32;NOT NULL;"`
}

func (RoleMenu) TableName() string {
	return "ag_role_menu"
}
