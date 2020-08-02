package global

import (
	"github.com/jinzhu/gorm"
)

type Model struct {
	gorm.Model
}

type OptLog struct {
	gorm.Model
	UserID   uint   `json:"user_id"`
	Url      string `json:"url"`
	Method   string `json:"method"`
	Body     string `json:"body"`
	ClientIP string `json:"client_ip"`
}

type AuthLog struct {
	gorm.Model
	Username   string `json:"username"`
	ClientIP   string `json:"client_ip"`
	AuthStatus uint   `json:"auth_status" gorm:"default:0"`
}
