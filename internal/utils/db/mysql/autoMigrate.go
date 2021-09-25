package mysql

import (
	"admin-go/internal/global"
	"admin-go/internal/groups"
	"admin-go/internal/menus"
	"admin-go/internal/roles"
	"admin-go/internal/users"
	"admin-go/internal/utils/db"
	_ "embed"
	"fmt"
	"strings"
)

var (
	//go:embed  init.sql
	initSQLs string
)

func AutoMigrate() {
	init := false
	if !db.Mysql.Migrator().HasTable("casbin_rule") {
		init = true
	}
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

	if init {
		initSQL()
	}
}

func initSQL() {
	sqls := strings.Split(initSQLs, ";")
	for i := range sqls {
		sql := strings.Replace(sqls[i], "\n", "", -1)
		if sql != "" {
			initErr := db.Mysql.Exec(sql).Error
			if initErr != nil {
				fmt.Println(initErr)
			}
		}
	}
}
