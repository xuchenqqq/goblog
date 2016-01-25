package controllers

import (
	"github.com/astaxie/beego"
	"github.com/smalltree0/goblog/RS"
)

type AuthController struct {
	beego.Controller
}

func (this *AuthController) Get() {
	this.TplNames = "login.html"
}

func (this *AuthController) Post() {
	resp := NewResponse()
	username := this.GetString("username")
	password := this.GetString("password")

	if username == "" || password == "" {
		resp.Status = RS.RS_params_error
		resp.Tips(WARNING, RS.RS_params_error)
		resp.WriteJson(this.Ctx.ResponseWriter)
		return
	}
	// if code := db.UMgr.LoginUser(username, password); code == RS.RS_user_inexistence {
	// 	resp.Status = code
	// 	resp.Tips(WARNING, code)
	// } else if code == RS.RS_password_error {
	// 	resp.Status = code
	// 	resp.Tips(WARNING, code)
	// } else {
	// 	sess, _ := gloablSessions.SessionStart(w, r)
	// 	sess.Set(SESSIONNAME, username)
	// 	resp.Data = "/"
	// }
	resp.WriteJson(this.Ctx.ResponseWriter)
}
