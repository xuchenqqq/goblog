package controllers

import (
	"fmt"
	"strconv"

	"github.com/deepzz0/goblog/models"
	// "github.com/deepzz0/go-common/log"
)

type TagController struct {
	Common
}

func (this *TagController) Get() {
	this.Layout = "homelayout.html"
	this.TplName = "groupTemplate.html"
	this.Leftbar("")
	this.ListTopic()
}

func (this *TagController) ListTopic() {
	tagName := this.Ctx.Input.Param(":tag")
	tag := models.Blogger.GetTagByID(tagName)
	this.Data["Name"] = "无效TAG"
	if tag != nil {
		this.Data["Name"] = tag.ID
		page := 1
		tagName := this.Ctx.Input.Param(":tag")
		pageStr := this.Ctx.Input.Param(":page")
		if temp, err := strconv.Atoi(pageStr); err == nil {
			page = temp
		}
		topics, remainpage := models.TMgr.GetTopicsByTag(tagName, page)
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
				this.Data["UrlOlder"] = this.domain + "/tag/" + tagName + fmt.Sprintf("/p/%d", page-1)
			}
			if remainpage == 0 {
				this.Data["ClassNewer"] = "disabled"
				this.Data["UrlNewer"] = "#"
			} else {
				this.Data["ClassNewer"] = ""
				this.Data["UrlNewer"] = this.domain + "/tag/" + tagName + fmt.Sprintf("/p/%d", page+1)
			}
			this.Data["ListTopics"] = topics
		}
	}
	this.Data["Title"] = tagName + " - " + models.Blogger.BlogName
	this.Data["Description"] = fmt.Sprintf("标签,%s,%s,blog", models.Blogger.Introduce, models.Blogger.UserName)
	this.Data["Keywords"] = fmt.Sprintf("标签,tag,%s,%s", models.Blogger.Introduce, models.Blogger.UserName)
}
