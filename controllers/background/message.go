package background

import (
	"github.com/deepzz0/goblog/models"
)

type MessageController struct {
	Common
}

func (this *MessageController) Get() {
	this.Layout = "manage/adminlayout.html"
	this.TplName = "manage/message.html"
	this.Data["Title"] = "留言管理 - " + models.Blogger.BlogName
	this.LeftBar("message")
	this.Content()
}

func (this *MessageController) Content() {
	this.Data["ID"] = 99999
	this.Data["URL"] = this.domain + "/message"
	this.Data["Title"] = "给我留言"
}
