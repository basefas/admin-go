package groups

import (
	"go-admin/internal/global"
	"time"
)

type Group struct {
	global.Model
	GroupName string `json:"name" gorm:"NOT NULL"`
}

func (Group) TableName() string {
	return "group_"
}

type CreateGroup struct {
	GroupName string `json:"name" binding:"required"`
	RoleID    uint   `json:"role_id" binding:"required"`
}

type UpdateGroup struct {
	GroupName string `json:"name"`
	RoleID    uint   `json:"role_id"`
}

type GetGroupInfo struct {
	ID        uint      `json:"id"`
	GroupName string    `json:"name"`
	RoleID    uint      `json:"role_id"`
	RoleName  string    `json:"role_name"`
	CreatedAt time.Time `json:"create_time"`
	UpdatedAt time.Time `json:"update_time"`
}

type GroupRole struct {
	global.Model
	GroupID uint `json:"group_id" gorm:"NOT NULL"`
	RoleID  uint `json:"role_id" gorm:"NOT NULL"`
}

type UserGroup struct {
	UserID  uint `json:"user_id"`
	GroupID uint `json:"group_id"`
}

type User struct {
	UserID   uint   `json:"id"`
	Username string `json:"name"`
}
