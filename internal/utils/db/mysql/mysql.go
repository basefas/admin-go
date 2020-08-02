package mysql

import (
	"fmt"
	"go-admin/internal/utils/db"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/viper"
)

func Init() {
	var err error
	user := viper.GetString("db.mysql.user")
	password := viper.GetString("db.mysql.password")
	host := viper.GetString("db.mysql.host")
	port := viper.GetString("db.mysql.port")
	name := viper.GetString("db.mysql.name")

	db.Mysql, err = gorm.Open("mysql",
		fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			user, password, host, port, name))
	fmt.Printf("Mysql: %s:%s Connection successful.\n", host, port)
	if err != nil {
		fmt.Println(err)
	}
	db.Mysql.SingularTable(true)

	//Migrate your schema if you need
	AutoMigrate()
}
