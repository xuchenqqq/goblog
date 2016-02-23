package controllers

import (
	"github.com/deepzz0/goblog/RS"
	"github.com/deepzz0/goblog/helper"
	"github.com/deepzz0/goblog/models"
)

type AuthController struct {
	BaseController
}

func (this *AuthController) Get() {
	if logout := this.GetString("logout"); logout == "now" {
		this.DelSession(sessionname)
	}
	this.TplName = "login.html"
	this.Data["Name"] = models.Blogger.BlogName
	this.Data["Url"] = this.domain
}

func (this *AuthController) Post() {
	resp := helper.NewResponse()
	username := this.GetString("username")
	password := this.GetString("password")

	if username == "" || password == "" {
		resp.Status = RS.RS_params_error
		resp.Tips(helper.WARNING, RS.RS_params_error)
		resp.WriteJson(this.Ctx.ResponseWriter)
		return
	}
	if code := models.UMgr.LoginUser(username, password); code == RS.RS_user_inexistence {
		resp.Status = code
		resp.Tips(helper.WARNING, code)
	} else if code == RS.RS_password_error {
		resp.Status = code
		resp.Tips(helper.WARNING, code)
	} else {
		models.Blogger.LoginIp = this.Ctx.Request.RemoteAddr
		this.SetSession(sessionname, username)
		resp.Data = "/admin/data"
	}
	resp.WriteJson(this.Ctx.ResponseWriter)
}
