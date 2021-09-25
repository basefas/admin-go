package global

import (
	"time"

	"gorm.io/gorm"
)

type Model struct {
	ID        uint64 `gorm:"primaryKey;type:uint;size:32;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	//gorm.Model
}

type OptLog struct {
	Model
	UserID   uint64 `json:"user_id" gorm:"type:uint;size:32;"`
	Url      string `json:"url"`
	Method   string `json:"method"`
	Body     string `json:"body"`
	ClientIP string `json:"client_ip"`
}

func (OptLog) TableName() string {
	return "ag_pot_log"
}

type AuthLog struct {
	Model
	Username   string `json:"username"`
	ClientIP   string `json:"client_ip"`
	AuthStatus uint64 `json:"auth_status" gorm:"default:0;type:uint;size:32;"`
}

func (AuthLog) TableName() string {
	return "ag_auth_log"
}
