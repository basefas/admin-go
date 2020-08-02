package mysql

import (
	"go-admin/internal/global"
	"go-admin/internal/groups"
	"go-admin/internal/menus"
	"go-admin/internal/roles"
	"go-admin/internal/users"
	"go-admin/internal/utils/db"
)

func AutoMigrate() {
	db.Mysql.AutoMigrate(&users.User{})
	db.Mysql.AutoMigrate(&groups.Group{})
	db.Mysql.AutoMigrate(&roles.Role{})
	db.Mysql.AutoMigrate(&users.UserGroup{})
	db.Mysql.AutoMigrate(&users.UserRole{})
	db.Mysql.AutoMigrate(&menus.Menu{})
	db.Mysql.AutoMigrate(&roles.RoleMenu{})
	db.Mysql.AutoMigrate(&global.AuthLog{})
	db.Mysql.AutoMigrate(&global.OptLog{})
}
