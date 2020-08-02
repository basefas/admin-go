package auth

import (
	"fmt"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v2"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var Casbin *casbin.Enforcer

func CheckPermission(c *gin.Context, e *casbin.Enforcer) bool {
	token := c.GetHeader("token")
	userID, _ := GetUID(token)
	path := c.Request.URL.Path
	method := c.Request.Method

	allowed, err := e.Enforce(fmt.Sprintf("user::%d", userID), path, method)
	if err != nil {
		fmt.Println(err)
	}
	return allowed
}

func Init() {
	m := model.NewModel()
	m.AddDef("r", "r", "sub, obj, act")
	m.AddDef("p", "p", "sub, obj, act")
	m.AddDef("g", "g", "_, _")
	m.AddDef("e", "e", "some(where (p.eft == allow))")
	m.AddDef("m", "m", `g(r.sub, p.sub) && keyMatch2(r.obj, p.obj) && regexMatch(r.act, p.act)`)

	p, _ := gormadapter.NewAdapter("mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
			viper.GetString("db.mysql.user"),
			viper.GetString("db.mysql.password"),
			viper.GetString("db.mysql.host"),
			viper.GetString("db.mysql.port"),
			viper.GetString("db.mysql.name")),
		true)
	Casbin, _ = casbin.NewEnforcer(m, p)

	err := Casbin.LoadPolicy()
	if err != nil {
		fmt.Println(err)
	}
}
