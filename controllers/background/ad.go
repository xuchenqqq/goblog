package background

import (
	"github.com/smalltree0/beego_goblog/models"
)

type ADController struct {
	BackgroundController
}

func (this *ADController) Get() {
	this.TplName = "manage/adminTemplate.html"
	this.Data["Title"] = "广告管理 - " + models.Blogger.BlogName
	this.LeftBar("ad")
	this.Content()
}

func (this *ADController) Content() {
	this.Data["Content"] = ""
}
