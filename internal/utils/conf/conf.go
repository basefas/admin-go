package conf

import (
	"fmt"

	"github.com/spf13/viper"
)

func Init() {
	viper.SetConfigName("app-config")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("[init] " + err.Error())
		panic("config init failed.")
	} else {
		fmt.Println("[init] Config Load Completed.")
	}
}
