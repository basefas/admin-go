package roles

import (
	"go-admin/internal/global"
	"time"
)

type Role struct {
	global.Model
	RoleName string `json:"name" gorm:"NOT NULL"`
}

func (Role) TableName() string {
	return "role_"
}

type CreateRole struct {
	RoleName string `json:"name" binding:"required"`
}

type UpdateRole struct {
	RoleName string `json:"name"`
}

type GetRoleInfo struct {
	ID        uint      `json:"id"`
	RoleName  string    `json:"name"`
	CreatedAt time.Time `json:"create_time"`
	UpdatedAt time.Time `json:"update_time"`
}

type RoleMenu struct {
	global.Model
	RoleID uint `json:"role_id" gorm:"NOT NULL"`
	MenuID uint `json:"menu_id" gorm:"NOT NULL"`
}
