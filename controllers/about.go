package controllers

import (
	"bytes"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/deepzz0/goblog/helper"
	"github.com/deepzz0/goblog/models"
)

type AboutController struct {
	BaseController
}

func (this *AboutController) Get() {
	this.TplName = "home.html"
	this.Leftbar("about")
	this.Content()
	this.Data["Title"] = "关于博主 - " + models.Blogger.BlogName
}

func (this *AboutController) Post() {
	resp := helper.NewResponse()
	if about := models.TMgr.GetTopic(1); about != nil {
		resp.Data = string(about.Content)
	} else {
		resp.Data = "博主真懒。"
	}
	resp.WriteJson(this.Ctx.ResponseWriter)
}

func (this *AboutController) Content() {
	aboutT := beego.BeeTemplates["about.html"]
	var buffer bytes.Buffer
	aboutT.Execute(&buffer, map[string]string{"Title": "关于博主", "Url": this.domain + "/about"})
	this.Data["Content"] = fmt.Sprintf("%s", buffer.Bytes())
}
