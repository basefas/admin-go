package roles

import (
	"admin-go/internal/global"
	"time"
)

type Role struct {
	global.Model
	RoleName string `json:"name" gorm:"NOT NULL"`
}

type CreateRole struct {
	RoleName string `json:"name" binding:"required"`
}

type UpdateRole struct {
	RoleName string `json:"name"`
}

type RoleInfo struct {
	ID        uint64    `json:"id"`
	RoleName  string    `json:"name"`
	CreatedAt time.Time `json:"create_time"`
	UpdatedAt time.Time `json:"update_time"`
}

type RoleMenu struct {
	global.Model
	RoleID uint64 `json:"role_id" gorm:"type:uint;size:32;NOT NULL;"`
	MenuID uint64 `json:"menu_id" gorm:"type:uint;size:32;NOT NULL;"`
}

type User struct {
	UserID   uint64 `json:"id"`
	Username string `json:"name"`
}
