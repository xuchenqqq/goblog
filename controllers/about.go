package controllers

import (
	"bytes"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/smalltree0/beego_goblog/cache"
	"github.com/smalltree0/beego_goblog/helper"
	"github.com/smalltree0/beego_goblog/models"
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
	resp.Data = cache.Cache.AboutContent // markdown
	resp.WriteJson(this.Ctx.ResponseWriter)
}

func (this *AboutController) Content() {
	aboutT := beego.BeeTemplates["about.html"]
	var buffer bytes.Buffer
	aboutT.Execute(&buffer, map[string]string{"Title": "关于博主", "Url": this.domain + "/about"})
	this.Data["Content"] = fmt.Sprintf("%s", buffer.Bytes())
}
