package mysql

import (
	"admin-go/internal/global"
	"admin-go/internal/groups"
	"admin-go/internal/menus"
	"admin-go/internal/roles"
	"admin-go/internal/users"
	"admin-go/internal/utils/db"
	"fmt"
)

func AutoMigrate() {
	err := db.Mysql.AutoMigrate(
		&users.User{},
		&users.UserGroup{},
		&users.UserRole{},
		&groups.Group{},
		&groups.GroupRole{},
		&roles.Role{},
		&roles.RoleMenu{},
		&menus.Menu{},
		&global.AuthLog{},
		&global.OptLog{})
	if err != nil {
		fmt.Println("[init] " + err.Error())
		panic("db migrate failed.")
	}
}
