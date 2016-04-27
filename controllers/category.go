package controllers

import (
	"fmt"
	"strconv"

	// "github.com/deepzz0/go-common/log"
	"github.com/deepzz0/goblog/models"
)

type CategoryController struct {
	Common
}

func (this *CategoryController) Get() {
	this.Layout = "homelayout.html"
	this.TplName = "groupTemplate.html"
	this.ListTopic()
}

func (this *CategoryController) ListTopic() {
	cat := this.Ctx.Input.Param(":cat")
	this.Leftbar(cat)
	category := models.Blogger.GetCategoryByID(cat)
	var name string = "暂无该分类"
	if category != nil && category.Extra != "" {
		name = category.Text
	}
	this.Data["Name"] = name
	this.Data["URL"] = fmt.Sprintf("%s/cat/%s", this.domain, category.ID)
	this.Data["Domain"] = this.domain
	pageStr := this.Ctx.Input.Param(":page")
	page := 1
	if temp, err := strconv.Atoi(pageStr); err == nil {
		page = temp
	}
	topics, remainpage := models.TMgr.GetTopicsByCatgory(cat, page)
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
			this.Data["UrlOlder"] = this.domain + "/cat/" + cat + fmt.Sprintf("/p/%d", page-1)
		}
		if remainpage == 0 {
			this.Data["ClassNewer"] = "disabled"
			this.Data["UrlNewer"] = "#"
		} else {
			this.Data["ClassNewer"] = ""
			this.Data["UrlNewer"] = this.domain + "/cat/" + cat + fmt.Sprintf("/p/%d", page+1)
		}
		this.Data["ListTopics"] = topics
	}
	this.Data["Title"] = name + " - " + models.Blogger.BlogName
	this.Data["Description"] = fmt.Sprintf("%s的个人博客,%s,%s,blog", models.Blogger.UserName, models.Blogger.Introduce, category.Text)
	this.Data["KeyWords"] = fmt.Sprintf("博客分类,%s,%s,%s", category.Text, models.Blogger.Introduce, models.Blogger.UserName)
}
