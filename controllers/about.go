package controllers

import (
	"github.com/deepzz0/goblog/models"
)

type AboutController struct {
	Common
}

func (this *AboutController) Get() {
	this.Layout = "homelayout.html"
	this.TplName = "aboutTemplate.html"
	this.Data["Title"] = "关于博主 - " + models.Blogger.BlogName
	this.Leftbar("about")
	this.Content()

}

func (this *AboutController) Content() {
	this.Data["Title"] = "关于博主"
	this.Data["URL"] = this.domain + "/about"
	if about := models.TMgr.GetTopic(1); about != nil {
		this.Data["Content"] = string(about.Content)
	} else {
		this.Data["Content"] = "博主真懒。"
	}
}
