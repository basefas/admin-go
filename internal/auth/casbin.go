package auth

import (
	"admin-go/internal/utils/log"
	"fmt"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var Casbin *casbin.Enforcer

func Init() {
	m := model.NewModel()
	m.AddDef("r", "r", "sub, obj, act")
	m.AddDef("p", "p", "sub, obj, act")
	m.AddDef("g", "g", "_, _")
	m.AddDef("e", "e", "some(where (p.eft == allow))")
	m.AddDef("m", "m", `g(r.sub, p.sub) && keyMatch2(r.obj, p.obj) && regexMatch(r.act, p.act)`)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		viper.GetString("db.mysql.user"),
		viper.GetString("db.mysql.password"),
		viper.GetString("db.mysql.host"),
		viper.GetString("db.mysql.port"),
		viper.GetString("db.mysql.name"))
	e, err := gormadapter.NewAdapter("mysql", dsn, true)
	if err != nil {
		fmt.Println("[init] " + err.Error())
		panic("auth init failed.")
	}
	Casbin, _ = casbin.NewEnforcer(m, e)

	errLoadPolicy := Casbin.LoadPolicy()
	if errLoadPolicy != nil {
		fmt.Println("[init] " + errLoadPolicy.Error())
		panic("auth init failed.")
	}
}

func CheckPermission(c *gin.Context, e *casbin.Enforcer) bool {
	token := c.GetHeader("token")
	userID, _ := GetUID(token)
	path := c.Request.URL.Path
	method := c.Request.Method

	allowed, err := e.Enforce(fmt.Sprintf("user::%d", userID), path, method)
	if err != nil {
		log.Warn(err.Error())
	}
	return allowed
}
