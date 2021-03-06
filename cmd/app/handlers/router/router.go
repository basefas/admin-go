package router

import (
	"admin-go/cmd/app/handlers/base"
	v1 "admin-go/cmd/app/handlers/v1"
	"admin-go/internal/auth"
	middleware "admin-go/internal/mid"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func Init() {
	r := setupRouter()
	port := fmt.Sprintf(":%d", viper.GetInt64("app.port"))
	err := r.Run(port)

	if err != nil {
		fmt.Println("[init] " + err.Error())
		panic("router init failed.")
	}
}

func setupRouter() *gin.Engine {
	runMode := viper.GetString("app.runMode")
	if runMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	r := gin.New()
	if runMode == "release" {
		r.Use(middleware.GinLogger())
		r.Use(middleware.GinRecovery(true))
	} else {
		r.Use(gin.Logger())
		r.Use(gin.Recovery())
	}
	r.Use(middleware.Cors())
	r.Use(middleware.Syslog())

	r.GET("/health", base.Health)
	api := r.Group("/api/v1")
	api.GET("/menus", v1.MenuList)
	api.POST("/login", v1.LogIn)

	user := api.Group("/user")
	user.Use(middleware.JWT())
	{
		user.GET("/:id", v1.UserGet)
		user.PUT("/:id", v1.UserUpdate)
	}

	users := api.Group("/users")
	users.Use(middleware.JWT())
	users.Use(middleware.Casbin(auth.Casbin))
	{
		users.POST("", v1.UsersCreate)
		users.GET("/:id", v1.UsersGet)
		users.PUT("/:id", v1.UsersUpdate)
		users.DELETE("/:id", v1.UsersDelete)
		users.GET("", v1.UsersList)
	}

	groups := api.Group("/groups")
	groups.Use(middleware.JWT())
	groups.Use(middleware.Casbin(auth.Casbin))
	{
		groups.POST("", v1.GroupCreate)
		groups.GET("/:id", v1.GroupGet)
		groups.PUT("/:id", v1.GroupUpdate)
		groups.DELETE("/:id", v1.GroupDelete)
		groups.GET("", v1.GroupList)
	}

	roles := api.Group("/roles")
	roles.Use(middleware.JWT())
	roles.Use(middleware.Casbin(auth.Casbin))
	{
		roles.POST("", v1.RoleCreate)
		roles.GET("/:id", v1.RoleGet)
		roles.PUT("/:id", v1.RoleUpdate)
		roles.DELETE("/:id", v1.RoleDelete)
		roles.GET("", v1.RoleList)
		roles.GET("/:id/menus", v1.RoleMenusList)
		roles.PUT("/:id/menus", v1.RoleMenusUpdate)
	}

	menus := api.Group("/menus")
	menus.Use(middleware.JWT())
	menus.Use(middleware.Casbin(auth.Casbin))
	{
		menus.POST("", v1.MenuCreate)
		menus.GET("/:id", v1.MenuGet)
		menus.PUT("/:id", v1.MenuUpdate)
		menus.DELETE("/:id", v1.MenuDelete)
	}

	return r
}
