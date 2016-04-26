package controllers

import (
	"fmt"
	"strconv"

	// "github.com/deepzz0/go-common/log"
	"github.com/deepzz0/goblog/models"
)

type HomeController struct {
	Common
}

func (this *HomeController) Get() {
	this.Layout = "homelayout.html"
	this.TplName = "homeTemplate.html"
	this.Data["Title"] = fmt.Sprintf("%s - %s", models.Blogger.Introduce, models.Blogger.BlogName)
	this.Leftbar("homepage")
	this.Home()
}
func (this *HomeController) Home() {
	this.Data["Tags"] = models.Blogger.Tags
	this.Data["Blogrolls"] = models.Blogger.Blogrolls
	this.Data["Domain"] = this.domain
	// 文章列表
	page := 1
	pageStr := this.Ctx.Input.Param(":page")
	if temp, err := strconv.Atoi(pageStr); err == nil {
		page = temp
	}
	topics, remainpage := models.TMgr.GetTopicsByPage(page)
	// log.Debugf("page = %d，remainpage=%d	", page, remainpage)
	if remainpage == -1 {
		this.Data["ClassOlder"] = "disabled"
		this.Data["UrlOlder"] = "#"
		this.Data["ClassNewer"] = "disabled"
		this.Data["UrlNewer"] = "#"
	} else {
		if page == 1 {
			this.Data["ClassOlder"] = "disabled"
			this.Data["UrlOlder"] = "#"
		} else {
			this.Data["ClassOlder"] = ""
			this.Data["UrlOlder"] = this.domain + "/p/" + fmt.Sprint(page-1)
		}
		if remainpage == 0 {
			this.Data["ClassNewer"] = "disabled"
			this.Data["UrlNewer"] = "#"
		} else {
			this.Data["ClassNewer"] = ""
			this.Data["UrlNewer"] = this.domain + "/p/" + fmt.Sprint(page+1)
		}
		this.Data["ListTopics"] = topics
	}
}
