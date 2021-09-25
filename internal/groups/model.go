package groups

import (
	"admin-go/internal/global"
	"time"
)

type Group struct {
	global.Model
	Name string `json:"name" gorm:"NOT NULL"`
}

func (Group) TableName() string {
	return "ag_group"
}

type CreateGroup struct {
	Name   string `json:"name" binding:"required"`
	RoleID uint64 `json:"role_id" binding:"required"`
}

type UpdateGroup struct {
	Name   string `json:"name"`
	RoleID uint64 `json:"role_id"`
}

type GroupInfo struct {
	ID        uint64    `json:"id"`
	Name      string    `json:"name"`
	RoleID    uint64    `json:"role_id"`
	RoleName  string    `json:"role_name"`
	CreatedAt time.Time `json:"create_time"`
	UpdatedAt time.Time `json:"update_time"`
}

type GroupRole struct {
	global.Model
	GroupID uint64 `json:"group_id" gorm:"type:uint;size:32;NOT NULL;"`
	RoleID  uint64 `json:"role_id" gorm:"type:uint;size:32;NOT NULL;"`
}

func (GroupRole) TableName() string {
	return "ag_group_role"
}

type UserGroup struct {
	UserID  uint64 `json:"user_id"`
	GroupID uint64 `json:"group_id"`
}

type User struct {
	UserID   uint64 `json:"id"`
	Username string `json:"name"`
}
