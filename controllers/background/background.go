package background

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/deepzz/beego_goblog/RS"
	"github.com/deepzz/beego_goblog/cache"
	"github.com/deepzz/beego_goblog/helper"
	// "github.com/deepzz/beego_goblog/models"
	// "github.com/deepzz/com/log"
)

var sessionname = beego.AppConfig.String("sessionname")

type BackgroundController struct {
	beego.Controller
	index  string
	domain string
	url    string
}

func (this *BackgroundController) Prepare() {
	this.url = this.Ctx.Request.URL.String()
	this.domain = beego.AppConfig.String("mydomain")
	if beego.BConfig.RunMode == beego.DEV {
		this.domain = this.domain + ":" + beego.AppConfig.String("httpport")
	}
}
func (this *BackgroundController) LeftBar(index string) {
	var html string
	for _, node := range cache.Cache.BackgroundLeftBars {
		if node.ID != "" {
			if node.ID == index {
				node.Node.Class = "active"
			} else {
				node.Node.Class = ""
			}
		}
		html += node.Node.String()
	}
	this.Data["LeftBar"] = html
}

// ----------------------------- 过滤登录 -----------------------------
var FilterUser = func(ctx *context.Context) {
	val, ok := ctx.Input.Session(sessionname).(string)
	if !ok || val == "" {
		if ctx.Request.Method == "GET" {
			ctx.Redirect(302, "/login")
		} else if ctx.Request.Method == "POST" {
			resp := helper.NewResponse()
			resp.Status = RS.RS_user_not_login
			resp.Data = "/login"
			resp.WriteJson(ctx.ResponseWriter)
		}
	}
}
