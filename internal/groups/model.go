package groups

import (
	"go-admin/internal/global"
	"time"
)

type Group struct {
	global.Model
	GroupName string `json:"name" gorm:"NOT NULL"`
	ParentID  uint   `json:"parent_id" gorm:"NOT NULL"`
}

func (Group) TableName() string {
	return "group_"
}

type CreateGroup struct {
	GroupName string `json:"name" binding:"required"`
	ParentID  uint   `json:"parent_id" binding:"required"`
}

type UpdateGroup struct {
	GroupName string `json:"name"`
	ParentID  uint   `json:"parent_id"`
}

type GetGroupInfo struct {
	ID        uint      `json:"id"`
	GroupName string    `json:"name"`
	HeadCount string    `json:"head_count"`
	ParentID  uint      `json:"parent_id"`
	CreatedAt time.Time `json:"create_time"`
	UpdatedAt time.Time `json:"update_time"`
}
