package users

import (
	"go-admin/internal/global"
	"time"
)

type User struct {
	global.Model
	Username string `json:"username" gorm:"NOT NULL"`
	Password string `json:"password" gorm:"NOT NULL"`
	Email    string `json:"email"`
	Status   uint   `json:"status"`
}

func (User) TableName() string {
	return "user_"
}

type CreateUser struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
	GroupID  uint   `json:"group_id"`
	RoleID   uint   `json:"role_id"`
	Status   uint   `json:"status"`
}

type UpdateUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	GroupID  uint   `json:"group_id"`
	RoleID   uint   `json:"role_id"`
	Status   uint   `json:"status"`
}

type UserInfo struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"create_time"`
	UpdatedAt time.Time `json:"update_time"`
	GroupID   uint      `json:"group_id"`
	GroupName string    `json:"group_name"`
	RoleID    uint      `json:"role_id"`
	RoleName  string    `json:"role_name"`
	Status    uint      `json:"status"`
}

type Login struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Token    string `json:"token"`
}

type UserGroup struct {
	global.Model
	UserID  uint `json:"user_id" gorm:"NOT NULL"`
	GroupID uint `json:"group_id" gorm:"NOT NULL"`
}

type UserRole struct {
	global.Model
	UserID uint `json:"user_id" gorm:"NOT NULL"`
	RoleID uint `json:"role_id" gorm:"NOT NULL"`
}
