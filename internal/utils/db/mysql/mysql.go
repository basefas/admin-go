package mysql

import (
	"admin-go/internal/utils/db"
	"fmt"
	"strings"

	"gorm.io/driver/mysql"

	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type DBMySQLError struct {
	Number  uint16
	Message string
}

func Init() {
	var err error
	user := viper.GetString("db.mysql.user")
	password := viper.GetString("db.mysql.password")
	host := viper.GetString("db.mysql.host")
	port := viper.GetUint64("db.mysql.port")
	name := viper.GetString("db.mysql.name")
	CreateDatabase(user, password, host, port, name)

	dsn := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, name)

	db.Mysql, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                      dsn,
		DefaultStringSize:        255,
		DisableDatetimePrecision: true,
	}), &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		if strings.Contains(err.Error(), "Error 1049") {
		} else {
			fmt.Println("[init] " + err.Error())
			panic("db init failed.")
		}
	} else {
		fmt.Printf("[init] Mysql %s:%d Connection Successful.\n", host, port)
	}
	AutoMigrate()
}

func CreateDatabase(user, password, host string, port uint64, name string) {

	dsn := fmt.Sprintf("%s:%s@(%s:%d)/?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port)

	t, err := gorm.Open(mysql.Open(dsn), nil)
	if err != nil {
		fmt.Println("[init] ", err.Error())
	}

	createSQL := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4;", name)

	err = t.Exec(createSQL).Error
	if err != nil {
		fmt.Println("[init] " + err.Error())
	}
}
