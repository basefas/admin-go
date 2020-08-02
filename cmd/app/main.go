package main

import (
	"go-admin/cmd/app/handlers"
	"go-admin/internal/auth"
	"go-admin/internal/utils/conf"
	"go-admin/internal/utils/db/mysql"
)

func main() {
	run()
}

func run() {
	conf.Init()
	mysql.Init()
	auth.Init()
	handlers.Init()
}
