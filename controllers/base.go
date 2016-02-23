package controllers

import (
	"time"

	"github.com/astaxie/beego"
	"github.com/deepzz0/go-common/log"
	"github.com/deepzz0/goblog/models"
)

var sessionname = beego.AppConfig.String("sessionname")

type BaseController struct {
	beego.Controller
	domain string
	url    string
}

func (this *BaseController) Prepare() {
	this.url = this.Ctx.Request.URL.String()
	this.domain = beego.AppConfig.String("mydomain")
	if beego.BConfig.RunMode == beego.DEV {
		this.domain = this.domain + ":" + beego.AppConfig.String("httpport")
	}
	log.Debugf("%s", this.domain)
}
func (this *BaseController) Leftbar(cat string) {
	var html string
	for _, n := range models.Blogger.Categories {
		if len(n.Node.Children) > 0 {
			if cat == n.ID {
				n.Node.Class = "active"
			} else {
				n.Node.Class = ""
			}
		}
		html += n.Node.String()
	}
	this.Data["Picture"] = models.Blogger.HeadIcon
	this.Data["BlogName"] = models.Blogger.BlogName
	this.Data["Introduce"] = models.Blogger.Introduce
	this.Data["Category"] = html
	html = ""
	for _, s := range models.Blogger.Socials {
		html += s.Node.String()
	}
	this.Data["Social"] = html
	this.Data["Url"] = this.domain
	this.Data["CopyTime"] = time.Now().Year()
}

type listOfTopic struct {
	ID       int32
	Title    string
	Url      string
	Time     string
	Preview  string
	Category string
	Tags     string
}
