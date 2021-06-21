package users

import (
	"admin-go/internal/global"
	"time"
)

type User struct {
	global.Model
	Username string `json:"username" gorm:"NOT NULL"`
	Password string `json:"password" gorm:"NOT NULL"`
	Email    string `json:"email"`
	Status   uint64 `json:"status" gorm:"type:uint;size:32;"`
}

type CreateUser struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
	GroupID  uint64 `json:"group_id"`
	RoleID   uint64 `json:"role_id"`
	Status   uint64 `json:"status"`
}

type UpdateUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	GroupID  uint64 `json:"group_id"`
	RoleID   uint64 `json:"role_id"`
	Status   uint64 `json:"status"`
}

type UserInfo struct {
	ID        uint64    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"create_time"`
	UpdatedAt time.Time `json:"update_time"`
	GroupID   uint64    `json:"group_id"`
	GroupName string    `json:"group_name"`
	RoleID    uint64    `json:"role_id"`
	RoleName  string    `json:"role_name"`
	Status    uint64    `json:"status"`
}

type Login struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginInfo struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Token    string `json:"token"`
}

type UserGroup struct {
	global.Model
	UserID  uint64 `json:"user_id" gorm:"type:uint;size:32;NOT NULL;"`
	GroupID uint64 `json:"group_id" gorm:"type:uint;size:32;NOT NULL;"`
}

type UserRole struct {
	global.Model
	UserID uint64 `json:"user_id" gorm:"type:uint;size:32;NOT NULL;"`
	RoleID uint64 `json:"role_id" gorm:"type:uint;size:32;NOT NULL;"`
}
