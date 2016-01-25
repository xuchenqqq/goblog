package routers

import (
	"github.com/astaxie/beego"
	"github.com/smalltree0/beego_goblog/controllers"
)

func init() {
	beego.SessionOn = true
	beego.SessionName = "SESSIONID"
	beego.SessionGCMaxLifetime = 3600
	beego.SessionCookieLifeTime = 3600

	beego.Router("/", &controllers.HomeController{})
	beego.Router("/login", &controllers.AuthController{})
}
