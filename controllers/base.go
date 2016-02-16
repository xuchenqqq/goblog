package controllers

import (
	// "encoding/json"
	// "fmt"
	"time"

	"github.com/astaxie/beego"
	// "github.com/smalltree0/beego_goblog/RS"
	// "github.com/smalltree0/beego_goblog/helper"
	"github.com/smalltree0/beego_goblog/models"
	// "github.com/smalltree0/com/log"
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
