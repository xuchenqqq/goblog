package controllers

import (
	"github.com/deepzz0/goblog/models"
)

type MessageController struct {
	Common
}

func (this *MessageController) Get() {
	this.Layout = "homelayout.html"
	this.TplName = "messageTemplate.html"
	this.Data["Title"] = "给我留言 - " + models.Blogger.BlogName
	this.Leftbar("message")
	this.Content()
}

func (this *MessageController) Content() {
	this.Data["Title"] = "给我留言"
	this.Data["ID"] = "99999"
	this.Data["URL"] = this.domain + "/message"
}
