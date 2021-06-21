package main

import (
	"admin-go/cmd/app/handlers/router"
	"admin-go/internal/auth"
	"admin-go/internal/utils/conf"
	"admin-go/internal/utils/db/mysql"
	"admin-go/internal/utils/log"
)

func main() {
	conf.Init()
	log.Init()
	mysql.Init()
	auth.Init()
	router.Init()
}
