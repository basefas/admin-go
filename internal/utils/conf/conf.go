package conf

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func Init() {
	loadConfig()
	watchConfig()
	fmt.Println("Config Load Completed.")
}

func loadConfig() {
	viper.SetConfigName("app-config")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}
}

func watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("Config changed: %s.\n", e.Name)
	})
}
